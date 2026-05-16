package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册内容管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/post/:id", authMiddleware, h.Detail)
	group.GET("/post", authMiddleware, h.List)
	group.POST("/post", authMiddleware, h.Create)
	group.PATCH("/post/:id", authMiddleware, h.Update)
	group.PUT("/post/:id", authMiddleware, h.Publish)
	group.DELETE("/post", authMiddleware, h.Delete)
}
