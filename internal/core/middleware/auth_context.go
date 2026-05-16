// Package middleware provides Gin middleware for CORS, JWT authentication, rate limiting, request logging, panic recovery, and context enrichment.
package middleware

// 本文件定义认证上下文结构以及读取辅助方法。

import "github.com/gin-gonic/gin"

// AuthContext 描述认证中间件从 JWT 中解析出的轻量身份信息。
type AuthContext struct {
	UserID    int64
	Username  string
	ProjectID string
}

// SetAuthContext 向上下文中写入认证信息。
func SetAuthContext(c *gin.Context, auth *AuthContext) {
	c.Set(ContextKeyAuth, auth)
}

// GetAuthContext 从上下文中读取认证信息。
func GetAuthContext(c *gin.Context) (*AuthContext, bool) {
	value, ok := c.Get(ContextKeyAuth)
	if !ok {
		return nil, false
	}
	auth, ok := value.(*AuthContext)
	return auth, ok && auth != nil
}

// GetUsername 从上下文中读取当前用户名。
func GetUsername(c *gin.Context) (string, bool) {
	auth, ok := GetAuthContext(c)
	if !ok || auth.Username == "" {
		return "", false
	}
	return auth.Username, true
}

// GetUserID 从上下文中读取当前用户 ID。
func GetUserID(c *gin.Context) (int64, bool) {
	auth, ok := GetAuthContext(c)
	if ok && auth.UserID > 0 {
		return auth.UserID, true
	}
	value, ok := c.Get(ContextKeyUserID)
	if !ok {
		return 0, false
	}
	userID, ok := value.(int64)
	return userID, ok
}
