package middleware

import (
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"gorm.io/gorm"

	apirepo "gotribe/internal/admin/api/repository"
	"gotribe/internal/core/database"
	"gotribe/internal/model"
)

// 操作日志channel
var OperationLogChan = make(chan *model.OperationLog, 256)
var apiDescCache = cache.New(10*time.Minute, 20*time.Minute)

// 定义静态资源路径前缀
var skipPaths = []string{
	"/static/",
	"/assets/",
	"/images/",
	"/favicon.ico",
	"/swagger/",
}

func OperationLogMiddleware(db *gorm.DB, enforcer *casbin.SyncedEnforcer, urlPathPrefix string, log *zap.SugaredLogger) gin.HandlerFunc {
	tx := database.NewTransactionManager(db)
	return func(c *gin.Context) {
		// 获取实际请求路径
		requestPath := c.Request.URL.Path

		// 如果请求的是 Swagger 相关路径，直接跳过
		if strings.HasPrefix(requestPath, "/swagger/") {
			c.Next()
			return
		}

		// 获取访问路径
		path := strings.TrimPrefix(c.FullPath(), "/"+urlPathPrefix)

		// 如果是空路径或静态资源，直接返回
		if shouldSkipLog(path) {
			c.Next()
			return
		}

		startTime := time.Now()
		c.Next()
		timeCost := time.Since(startTime).Milliseconds()

		username := getUsername(c)
		method := c.Request.Method

		// 获取接口描述
		apiDesc := getApiDescription(path, method, c, tx, enforcer, log)

		operationLog := &model.OperationLog{
			Username:  username,
			Ip:        c.ClientIP(),
			Method:    method,
			Path:      path,
			Desc:      apiDesc,
			Status:    c.Writer.Status(),
			StartTime: startTime,
			TimeCost:  timeCost,
		}

		// 异步写入日志
		select {
		case OperationLogChan <- operationLog:
		default:
			// 日志写入不应阻塞主请求链路，队列满时丢弃并记录告警。
			if requestPath != "/health" {
				repositoryLogDropWarn(log, path, method)
			}
		}
	}
}

// 判断是否需要跳过日志记录
func shouldSkipLog(path string) bool {
	if path == "" {
		return true
	}

	for _, prefix := range skipPaths {
		if strings.HasPrefix(path, prefix) || prefix == path {
			return true
		}
	}
	return false
}

// 获取用户名
func getUsername(c *gin.Context) string {
	ctxUser, exists := c.Get("user")
	if !exists {
		return "未登录"
	}

	user, ok := ctxUser.(model.Admin)
	if !ok {
		return "未登录"
	}

	return user.Username
}

// 获取API描述
func getApiDescription(path, method string, c *gin.Context, tx *database.TransactionManager, enforcer *casbin.SyncedEnforcer, log *zap.SugaredLogger) string {
	cacheKey := method + ":" + path
	if cachedDesc, ok := apiDescCache.Get(cacheKey); ok {
		return cachedDesc.(string)
	}

	apiRepository := apirepo.NewRepository(tx, enforcer)
	apiDesc, err := apiRepository.GetApiDescByPath(c.Request.Context(), path, method)
	if err != nil {
		if log != nil {
			log.Warnf("获取API描述失败: path=%s, method=%s, err=%v", path, method, err)
		}
		return ""
	}
	apiDescCache.Set(cacheKey, apiDesc, cache.DefaultExpiration)
	return apiDesc
}

func repositoryLogDropWarn(log *zap.SugaredLogger, path, method string) {
	if log == nil {
		return
	}
	// 使用固定文案降低告警噪音，避免在高压场景下再次放大日志量。
	if path == "" {
		log.Warn("operation log queue is full, dropping request log")
		return
	}
	log.Warnf("operation log queue is full, dropping request log for %s %s", method, path)
}
