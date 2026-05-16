package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册评论管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/comment", authMiddleware, h.List)
	group.PATCH("/comment/:id", authMiddleware, h.Update)
	group.DELETE("/comment", authMiddleware, h.Delete)
}
