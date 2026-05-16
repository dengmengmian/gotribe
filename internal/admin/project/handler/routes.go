package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册项目管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/project/:id", authMiddleware, h.Detail)
	group.GET("/project", authMiddleware, h.List)
	group.POST("/project", authMiddleware, h.Create)
	group.PATCH("/project/:id", authMiddleware, h.Update)
	group.DELETE("/project", authMiddleware, h.Delete)
}
