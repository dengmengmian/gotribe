package handler

import (
	"gotribe/internal/admin/index/service"
	"gotribe/internal/core/response"

	"github.com/gin-gonic/gin"
)

const defaultProjectID = "1"

// Handler 仪表盘处理器。
type Handler struct {
	indexService service.Service
}

// NewHandler 创建仪表盘处理器实例。
func NewHandler(indexService service.Service) *Handler {
	return &Handler{indexService: indexService}
}

// Dashboard 返回仪表盘全量数据。
func (h *Handler) Dashboard(c *gin.Context) {
	projectID := c.DefaultQuery("project_id", defaultProjectID)
	data, err := h.indexService.Dashboard(c.Request.Context(), projectID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{
		"indexDate": data,
	})
}

// CacheClear 清空 Redis 缓存。
func (h *Handler) CacheClear(c *gin.Context) {
	if err := h.indexService.CacheClear(c.Request.Context()); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"message": "缓存已清空"})
}
