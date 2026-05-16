package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册管理员管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/admin/info", authMiddleware, h.Me)
	group.GET("/admin/list", authMiddleware, h.List)
	group.PUT("/admin/changePwd", authMiddleware, h.UpdatePassword)
	group.POST("/admin/create", authMiddleware, h.Create)
	group.PATCH("/admin/update/:userID", authMiddleware, h.Update)
	group.GET("/admin/detail/:id", authMiddleware, h.Detail)
	group.DELETE("/admin/delete/batch", authMiddleware, h.Delete)
}
