package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册专栏管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/column", authMiddleware, h.List)
	group.GET("/column/:id", authMiddleware, h.Detail)
	group.POST("/column", authMiddleware, h.Create)
	group.PATCH("/column/:id", authMiddleware, h.Update)
	group.DELETE("/column", authMiddleware, h.Delete)
}
