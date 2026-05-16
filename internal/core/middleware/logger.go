package middleware

// 本文件实现 HTTP 请求访问日志中间件。

import (
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"gotribe/internal/core/config"
	"gotribe/internal/core/errs"
	applog "gotribe/internal/core/logger"
)

// Logger 创建 HTTP 请求日志中间件。
func Logger(app config.AppConfig) gin.HandlerFunc {
	verbose := app.IsDevelopment()

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		if !verbose && status < 400 && len(c.Errors) == 0 {
			return
		}
		args := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"query", sanitizeQuery(c.Request.URL.RawQuery),
			"status", status,
			"latency_ms", time.Since(start).Milliseconds(),
			"ip", c.ClientIP(),
			"errors", len(c.Errors),
		}
		if len(c.Errors) > 0 {
			args = append(args, "error_messages", c.Errors.String())
			if appErr := latestAppError(c); appErr != nil {
				args = append(args,
					"error_code", string(appErr.Code),
					"error_message", appErr.Message,
				)
				if len(appErr.Details) > 0 {
					args = append(args, "error_details", appErr.Details)
				}
				if appErr.Err != nil {
					args = append(args, "error_cause", appErr.Err.Error())
				}
			}
		}

		switch {
		case status >= 500:
			applog.Error(c.Request.Context(), "http request failed", args...)
		case status >= 400:
			applog.Warn(c.Request.Context(), "http request rejected", args...)
		default:
			applog.Info(c.Request.Context(), "http request completed", args...)
		}
	}
}

// sanitizeQuery 对敏感查询参数进行脱敏，避免将密码等值写入日志。
func sanitizeQuery(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return "[invalid_query_redacted]"
	}

	for _, key := range []string{"password", "token", "access_token", "refresh_token", "authorization", "auth_token", "code"} {
		if _, exists := values[key]; exists {
			values.Set(key, "REDACTED")
		}
	}
	return values.Encode()
}

// latestAppError 返回当前请求链路里最近一次记录的统一业务错误。
func latestAppError(c *gin.Context) *errs.AppError {
	if c == nil || len(c.Errors) == 0 {
		return nil
	}
	return errs.As(c.Errors.Last().Err)
}
