package handler

// 本文件负责注册当前用户模块的路由。

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册当前模块对外暴露的路由。
func (h *Handler) RegisterRoutes(secured *gin.RouterGroup, currentUser *gin.RouterGroup, apiLimiter gin.HandlerFunc, authLimiter gin.HandlerFunc) {
	secured.PATCH("/me", apiLimiter, h.UpdateMe)
	secured.POST("/me/password", authLimiter, h.ChangePassword)

	currentUser.GET("/me", apiLimiter, h.GetMe)
}
