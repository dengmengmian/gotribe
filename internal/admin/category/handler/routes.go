package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册分类管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/category/tree", authMiddleware, h.Tree)
	group.GET("/category", authMiddleware, h.List)
	group.POST("/category", authMiddleware, h.Create)
	group.PATCH("/category/:id", authMiddleware, h.Update)
	group.DELETE("/category", authMiddleware, h.Delete)
	group.GET("/category/:id", authMiddleware, h.Detail)
}
