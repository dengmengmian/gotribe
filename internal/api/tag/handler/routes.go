package handler

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册标签路由到指定路由组。
func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/tags", h.List)
}
