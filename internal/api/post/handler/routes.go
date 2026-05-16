package handler

// 本文件负责注册帖子模块的路由。

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册当前模块对外暴露的路由。
func (h *Handler) RegisterRoutes(public *gin.RouterGroup, apiLimiter gin.HandlerFunc) {
	public.GET("/posts", apiLimiter, h.List)
	public.GET("/posts/:postID", apiLimiter, h.Detail)
}
