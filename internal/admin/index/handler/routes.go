package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册首页管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/index/dashboard", authMiddleware, h.Dashboard)
	group.POST("/index/cache/clear", authMiddleware, h.CacheClear)
}
