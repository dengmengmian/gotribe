package handler

// 本文件负责注册 example 模块路由。

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册 example 模块对外暴露的路由。
func (h *Handler) RegisterRoutes(currentUser *gin.RouterGroup, apiLimiter gin.HandlerFunc) {
	currentUser.POST("/examples", apiLimiter, h.Create)
	currentUser.GET("/examples", apiLimiter, h.List)
	currentUser.GET("/examples/:exampleID", apiLimiter, h.Detail)
	currentUser.PATCH("/examples/:exampleID", apiLimiter, h.Update)
	currentUser.DELETE("/examples/:exampleID", apiLimiter, h.Delete)
}
