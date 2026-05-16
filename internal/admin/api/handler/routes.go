package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册接口管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/api/list", authMiddleware, h.List)
	group.GET("/api/tree", authMiddleware, h.Tree)
	group.POST("/api/create", authMiddleware, h.Create)
	group.PATCH("/api/update/:apiID", authMiddleware, h.Update)
	group.DELETE("/api/delete/batch", authMiddleware, h.Delete)
}
