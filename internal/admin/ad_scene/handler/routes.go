package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册推广场景管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/ad/scene", authMiddleware, h.List)
	group.GET("/ad/scene/:id", authMiddleware, h.Detail)
	group.POST("/ad/scene", authMiddleware, h.Create)
	group.PATCH("/ad/scene/:id", authMiddleware, h.Update)
	group.DELETE("/ad/scene", authMiddleware, h.Delete)
}
