package handler

// 本文件负责注册认证模块的路由。

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册当前模块对外暴露的路由。
func (h *Handler) RegisterRoutes(authRoutes *gin.RouterGroup, authLimiter gin.HandlerFunc, jwt gin.HandlerFunc) {
	authRoutes.POST("/login", authLimiter, h.Login)
	authRoutes.POST("/refresh", authLimiter, h.Refresh)
	authRoutes.POST("/logout", jwt, h.Logout)
}
