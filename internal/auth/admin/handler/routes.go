package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册认证管理路由
// 注意：实际的认证逻辑由 JWT 中间件处理，路由在 base_routes.go 中定义
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.POST("/base/login", authMiddleware, h.Login)
	group.POST("/base/logout", authMiddleware, h.Logout)
	group.POST("/base/refreshToken", authMiddleware, h.RefreshToken)
}
