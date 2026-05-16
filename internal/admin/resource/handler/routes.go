package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册资源管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/resource/:id", authMiddleware, h.Detail)
	group.GET("/resource", authMiddleware, h.List)
	group.PATCH("/resource/:id", authMiddleware, h.Update)
	group.POST("/resource/upload", authMiddleware, h.Upload)
	group.DELETE("/resource", authMiddleware, h.Delete)
}
