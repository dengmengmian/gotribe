package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册积分管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/point", authMiddleware, h.List)
	group.DELETE("/point", authMiddleware, h.Delete)
	group.GET("/point/:id", authMiddleware, h.Detail)
	group.PATCH("/point/:id", authMiddleware, h.Update)
	group.POST("/point", authMiddleware, h.Create)
}
