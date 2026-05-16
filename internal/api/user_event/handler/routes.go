// Package handler implements HTTP request handlers for user event reporting.
package handler

// 本文件负责注册用户行为模块的路由。

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册当前模块对外暴露的路由。
func (h *Handler) RegisterRoutes(secured *gin.RouterGroup, eventLimiter gin.HandlerFunc) {
	secured.POST("/user-events", eventLimiter, h.Create)
}
