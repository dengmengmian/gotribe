package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册广告管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/ad", authMiddleware, h.List)
	group.GET("/ad/:id", authMiddleware, h.Detail)
	group.POST("/ad", authMiddleware, h.Create)
	group.PATCH("/ad/:id", authMiddleware, h.Update)
	group.DELETE("/ad", authMiddleware, h.Delete)
}
