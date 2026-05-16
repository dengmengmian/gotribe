// Package handler provides HTTP handlers for the tag module.
package handler

import (
	"net/http"

	"gotribe/internal/core/response"
	"gotribe/internal/api/tag/dto"
	"gotribe/internal/api/tag/repository"

	"github.com/gin-gonic/gin"
)

// Handler 标签 HTTP 处理器。
type Handler struct {
	repo *repository.Repository
}

// NewHandler 创建标签处理器。
func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// List 返回所有启用标签。
func (h *Handler) List(c *gin.Context) {
	var req dto.ListQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.PerPage <= 0 {
		req.PerPage = 100
	}

	tags, err := h.repo.List(c.Request.Context(), req.Keyword, req.PerPage)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"tags": dto.ToTagListResponse(tags)})
}
