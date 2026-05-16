// Package handler provides HTTP handlers for the category module.
package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	categorydto "gotribe/internal/api/category/dto"
	"gotribe/internal/api/category/repository"
	"gotribe/internal/core/response"
)

// Handler 分类 HTTP 处理器。
type Handler struct {
	repo *repository.Repository
}

// NewHandler 创建分类处理器。
func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// ListByParentID 获取指定父分类下的子分类列表。
func (h *Handler) ListByParentID(c *gin.Context) {
	parentIDStr := c.Param("parent_id")
	parentID, err := strconv.ParseInt(parentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid parent_id"})
		return
	}

	categories, err := h.repo.ListByParentID(c.Request.Context(), parentID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, gin.H{"categories": categorydto.ToCategoryListResponse(categories)})
}
