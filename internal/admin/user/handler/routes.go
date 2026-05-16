package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册用户管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/user/:id", authMiddleware, h.Detail)
	group.GET("/user", authMiddleware, h.List)
	group.POST("/user", authMiddleware, h.Create)
	group.PATCH("/user/:id", authMiddleware, h.Update)
	group.DELETE("/user", authMiddleware, h.Delete)
	group.GET("/user/search", authMiddleware, h.Search)
}
