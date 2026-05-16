package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册菜单管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/menu/list", authMiddleware, h.List)
	group.GET("/menu/tree", authMiddleware, h.Tree)
	group.POST("/menu/create", authMiddleware, h.Create)
	group.PATCH("/menu/update/:menuID", authMiddleware, h.Update)
	group.DELETE("/menu/delete/batch", authMiddleware, h.Delete)
	group.GET("/menu/access/list/:userID", authMiddleware, h.GetUserMenusByUserID)
	group.GET("/menu/access/tree/:userID", authMiddleware, h.GetUserMenuTreeByUserID)
}
