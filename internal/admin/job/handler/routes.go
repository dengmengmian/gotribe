package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册任务管理路由
func (h *Handler) RegisterRoutes(group *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	group.GET("/jobs", authMiddleware, h.ListJobs)
	group.GET("/jobs/:name/status", authMiddleware, h.GetJobStatus)
	group.GET("/jobs/:name/history", authMiddleware, h.GetJobHistory)
	group.POST("/jobs/:name/enable", authMiddleware, h.EnableJob)
	group.POST("/jobs/:name/disable", authMiddleware, h.DisableJob)
}
