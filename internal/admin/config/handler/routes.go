package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册配置管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/config", authMiddleware, h.List)
	group.GET("/config/alias/:alias", authMiddleware, h.DetailByAlias)
	group.GET("/config/:id", authMiddleware, h.Detail)
	group.POST("/config", authMiddleware, h.Create)
	group.PATCH("/config/:id", authMiddleware, h.Update)
	group.DELETE("/config", authMiddleware, h.Delete)
}
