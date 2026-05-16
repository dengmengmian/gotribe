package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册角色管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/role/list", authMiddleware, h.List)
	group.POST("/role/create", authMiddleware, h.Create)
	group.PATCH("/role/update/:roleID", authMiddleware, h.Update)
	group.GET("/role/menus/get/:roleID", authMiddleware, h.GetRoleMenusByID)
	group.PATCH("/role/menus/update/:roleID", authMiddleware, h.UpdateRoleMenusByID)
	group.GET("/role/apis/get/:roleID", authMiddleware, h.GetRoleApisByID)
	group.PATCH("/role/apis/update/:roleID", authMiddleware, h.UpdateRoleApisByID)
	group.DELETE("/role/delete/batch", authMiddleware, h.Delete)
}
