package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册反馈管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/feedback", authMiddleware, h.List)
	group.DELETE("/feedback", authMiddleware, h.Delete)
}
