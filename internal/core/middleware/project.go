package middleware

// 本文件实现项目隔离信息注入中间件。

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/logger"
	"gotribe/internal/core/response"
)

// ProjectID 创建项目 ID 透传中间件。
func ProjectID(defaultProjectID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := strings.TrimSpace(c.GetHeader("X-Project-ID"))
		if projectID == "" {
			projectID = defaultProjectID
		}
		if projectID == "" {
			response.Error(c, errs.BadRequest("X-Project-ID is required", nil))
			c.Abort()
			return
		}
		c.Set(ContextKeyProjectID, projectID)
		c.Request = c.Request.WithContext(logger.WithProjectID(c.Request.Context(), projectID))
		c.Next()
	}
}

// GetProjectID 从上下文中读取当前项目 ID。
func GetProjectID(c *gin.Context) string {
	return c.GetString(ContextKeyProjectID)
}
