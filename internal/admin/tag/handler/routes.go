package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册标签管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/tag", authMiddleware, h.List)
	group.GET("/tag/:id", authMiddleware, h.Detail)
	group.POST("/tag", authMiddleware, h.Create)
	group.PATCH("/tag/:id", authMiddleware, h.Update)
	group.DELETE("/tag", authMiddleware, h.Delete)
}
