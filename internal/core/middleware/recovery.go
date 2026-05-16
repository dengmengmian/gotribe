package middleware

// 本文件实现统一的 panic 恢复中间件。

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gotribe/internal/core/errs"
	applog "gotribe/internal/core/logger"
	"gotribe/internal/core/response"
)

// Recovery 创建 panic 恢复中间件。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				applog.Error(c.Request.Context(), "panic recovered", "panic", fmt.Sprint(rec))
				response.Error(c, errs.Internal("internal server error", nil))
				c.Abort()
			}
		}()
		c.Next()
	}
}
