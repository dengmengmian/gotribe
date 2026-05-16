package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册操作日志管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/operation/list", authMiddleware, h.List)
	group.DELETE("/operation/delete/batch", authMiddleware, h.Delete)
}
