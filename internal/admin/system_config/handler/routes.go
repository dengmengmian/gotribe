package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册系统配置管理路由（GET /base/config 在 routes.go 中公开注册）
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.PATCH("/system", authMiddleware, h.Update)
}
