// Package observability provides metrics collection, tracing middleware, and Prometheus-compatible metric export.
package observability

// 本文件负责统一初始化 metrics、tracing 和 HTTP 可观测性中间件。

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"gotribe/internal/core/config"
	"gotribe/internal/core/logger"
)

// ShutdownFn 表示 observability 资源的关闭函数。
type ShutdownFn func(context.Context) error

type runtimeState struct {
	serviceName    string
	metricsEnabled bool
	tracingEnabled bool
}

type requestMetric struct {
	Method string
	Route  string
	Status string
}

type inflightMetric struct {
	Method string
	Route  string
}

type histogram struct {
	Bounds []float64
	Counts []uint64
	Sum    float64
	Count  uint64
}

type metricsState struct {
	mu             sync.RWMutex
	totals         map[requestMetric]uint64
	inflight       map[inflightMetric]int64
	durationByCode map[requestMetric]*histogram
}

var (
	state = runtimeState{
		serviceName:    "gotribe",
		metricsEnabled: true,
		tracingEnabled: true,
	}
	metrics = metricsState{
		totals:         map[requestMetric]uint64{},
		inflight:       map[inflightMetric]int64{},
		durationByCode: map[requestMetric]*histogram{},
	}
	httpDurationBounds = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5}
)

// Init 初始化全局可观测性能力。
// 当前实现不依赖外部后端，即使开启 tracing 也不会影响服务启动。
func Init(cfg config.Config) ShutdownFn {
	state.serviceName = strings.TrimSpace(cfg.App.Name)
	if state.serviceName == "" {
		state.serviceName = "gotribe"
	}
	state.metricsEnabled = cfg.Observability.MetricsEnabled
	state.tracingEnabled = cfg.Observability.TracingEnabled

	return func(context.Context) error { return nil }
}

// MetricsMiddleware 记录 HTTP 请求计数、耗时和并发数。
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !state.metricsEnabled {
			c.Next()
			return
		}

		start := time.Now()
		method := c.Request.Method
		rawRoute := c.Request.URL.Path
		inflightKey := inflightMetric{Method: method, Route: rawRoute}

		metrics.mu.Lock()
		metrics.inflight[inflightKey]++
		metrics.mu.Unlock()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = rawRoute
		}
		reqKey := requestMetric{Method: method, Route: route, Status: strconv.Itoa(c.Writer.Status())}

		metrics.mu.Lock()
		metrics.totals[reqKey]++
		h := metrics.durationByCode[reqKey]
		if h == nil {
			h = newHistogram(httpDurationBounds)
			metrics.durationByCode[reqKey] = h
		}
		h.observe(time.Since(start).Seconds())
		metrics.inflight[inflightKey]--
		if metrics.inflight[inflightKey] <= 0 {
			delete(metrics.inflight, inflightKey)
		}
		metrics.mu.Unlock()
	}
}

// TracingMiddleware 为每个请求生成 trace_id / span_id 并写入日志上下文与响应头。
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !state.tracingEnabled {
			c.Next()
			return
		}

		traceID := incomingTraceID(c.GetHeader("traceparent"))
		if traceID == "" {
			traceID = newTraceID()
		}
		spanID := newSpanID()

		ctx := logger.WithTrace(c.Request.Context(), traceID, spanID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Trace-ID", traceID)
		c.Writer.Header().Set("X-Span-ID", spanID)
		c.Writer.Header().Set("traceparent", fmt.Sprintf("00-%s-%s-01", traceID, spanID))
		c.Next()
	}
}

// MetricsHandler 返回 Prometheus 文本格式的指标导出端点。
func MetricsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		metrics.mu.RLock()
		defer metrics.mu.RUnlock()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		fmt.Fprintln(w, "# HELP gotribe_http_requests_total HTTP 请求总数。")
		fmt.Fprintln(w, "# TYPE gotribe_http_requests_total counter")
		for _, key := range sortedRequestMetrics(metrics.totals) {
			fmt.Fprintf(w, "gotribe_http_requests_total{method=%q,route=%q,status=%q} %d\n",
				key.Method, key.Route, key.Status, metrics.totals[key])
		}

		fmt.Fprintln(w, "# HELP gotribe_http_request_duration_seconds HTTP 请求耗时分布。")
		fmt.Fprintln(w, "# TYPE gotribe_http_request_duration_seconds histogram")
		for _, key := range sortedRequestMetrics(metrics.durationByCode) {
			h := metrics.durationByCode[key]
			var cumulative uint64
			for idx, bound := range h.Bounds {
				cumulative += h.Counts[idx]
				fmt.Fprintf(w, "gotribe_http_request_duration_seconds_bucket{method=%q,route=%q,status=%q,le=%q} %d\n",
					key.Method, key.Route, key.Status, trimFloat(bound), cumulative)
			}
			fmt.Fprintf(w, "gotribe_http_request_duration_seconds_bucket{method=%q,route=%q,status=%q,le=%q} %d\n",
				key.Method, key.Route, key.Status, "+Inf", h.Count)
			fmt.Fprintf(w, "gotribe_http_request_duration_seconds_sum{method=%q,route=%q,status=%q} %s\n",
				key.Method, key.Route, key.Status, trimFloat(h.Sum))
			fmt.Fprintf(w, "gotribe_http_request_duration_seconds_count{method=%q,route=%q,status=%q} %d\n",
				key.Method, key.Route, key.Status, h.Count)
		}

		fmt.Fprintln(w, "# HELP gotribe_http_requests_in_flight 当前正在处理中的 HTTP 请求数。")
		fmt.Fprintln(w, "# TYPE gotribe_http_requests_in_flight gauge")
		for _, key := range sortedInflightMetrics(metrics.inflight) {
			fmt.Fprintf(w, "gotribe_http_requests_in_flight{method=%q,route=%q} %d\n",
				key.Method, key.Route, metrics.inflight[key])
		}
	})
}

// InstrumentRedis 为 Redis 预留观测接入点；当前实现不要求额外 exporter。
func InstrumentRedis(*redis.Client) error { return nil }

// InstrumentGORM 为 GORM 预留观测接入点；当前实现不要求额外 exporter。
func InstrumentGORM(*gorm.DB) error { return nil }

func newHistogram(bounds []float64) *histogram {
	return &histogram{
		Bounds: append([]float64(nil), bounds...),
		Counts: make([]uint64, len(bounds)),
	}
}

func (h *histogram) observe(value float64) {
	h.Count++
	h.Sum += value
	for idx, bound := range h.Bounds {
		if value <= bound {
			h.Counts[idx]++
			return
		}
	}
}

func incomingTraceID(traceparent string) string {
	parts := strings.Split(traceparent, "-")
	if len(parts) != 4 || len(parts[1]) != 32 {
		return ""
	}
	return strings.ToLower(parts[1])
}

func newTraceID() string {
	return newHex(16)
}

func newSpanID() string {
	return newHex(8)
}

func newHex(size int) string {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return strings.Repeat("0", size*2)
	}
	return hex.EncodeToString(buf)
}

func trimFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', -1, 64)
}

func sortedRequestMetrics[T any](source map[requestMetric]T) []requestMetric {
	keys := make([]requestMetric, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Route != keys[j].Route {
			return keys[i].Route < keys[j].Route
		}
		if keys[i].Method != keys[j].Method {
			return keys[i].Method < keys[j].Method
		}
		return keys[i].Status < keys[j].Status
	})
	return keys
}

func sortedInflightMetrics(source map[inflightMetric]int64) []inflightMetric {
	keys := make([]inflightMetric, 0, len(source))
	for key := range source {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Route != keys[j].Route {
			return keys[i].Route < keys[j].Route
		}
		return keys[i].Method < keys[j].Method
	})
	return keys
}
