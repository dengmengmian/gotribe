package handler

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册分类路由到指定路由组。
func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/categories/:parent_id", h.ListByParentID)
}
