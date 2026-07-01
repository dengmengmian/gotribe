// Package handler provides HTTP handlers for the tag module.
package handler

import (
	"gotribe/internal/api/tag/dto"
	tagservice "gotribe/internal/api/tag/service"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

// Handler 标签 HTTP 处理器。
type Handler struct {
	service *tagservice.Service
}

// NewHandler 创建标签处理器。
func NewHandler(service *tagservice.Service) *Handler {
	return &Handler{service: service}
}

// List 返回所有启用标签。
func (h *Handler) List(c *gin.Context) {
	var req dto.ListQuery
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	tags, err := h.service.List(c.Request.Context(), req.Keyword, req.PerPage)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"tags": dto.ToTagListResponse(tags)})
}
