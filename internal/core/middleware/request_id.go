// Package middleware provides Gin middleware for CORS, JWT authentication, rate limiting, request logging, panic recovery, and context enrichment.
package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"gotribe/internal/core/logger"
)

const (
	// ContextKeyRequestID is the Gin context key used to store the request ID.
	ContextKeyRequestID = "request_id"
	// ContextKeyProjectID is the Gin context key used to store the project ID.
	ContextKeyProjectID = "project_id"
	// ContextKeyUserID is the Gin context key used to store the user ID.
	ContextKeyUserID = "user_id"
	// ContextKeyUsername is the Gin context key used to store the username.
	ContextKeyUsername = "username"
	// ContextKeyAuth is the Gin context key used to store the parsed auth context.
	ContextKeyAuth = "auth_context"
	// ContextKeyCurrentUser is the Gin context key used to store the fully loaded current user.
	ContextKeyCurrentUser = "current_user"
)

// RequestID 创建请求 ID 注入中间件。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Set(ContextKeyRequestID, requestID)
		c.Request = c.Request.WithContext(logger.WithRequestID(c.Request.Context(), requestID))
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// newRequestID 生成新的请求追踪标识。
func newRequestID() string {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "req_fallback"
	}
	return "req_" + hex.EncodeToString(buf)
}
