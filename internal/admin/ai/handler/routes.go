package handler

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册通用 AI 路由。
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.POST("/ai/generate", authMiddleware, h.Generate)
}
