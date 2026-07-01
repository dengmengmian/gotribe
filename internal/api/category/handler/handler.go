// Package handler provides HTTP handlers for the category module.
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	categorydto "gotribe/internal/api/category/dto"
	categoryservice "gotribe/internal/api/category/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
)

// Handler 分类 HTTP 处理器。
type Handler struct {
	service *categoryservice.Service
}

// NewHandler 创建分类处理器。
func NewHandler(service *categoryservice.Service) *Handler {
	return &Handler{service: service}
}

// ListByParentID 获取指定父分类下的子分类列表。
func (h *Handler) ListByParentID(c *gin.Context) {
	parentID, err := strconv.ParseInt(c.Param("parent_id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("invalid parent_id", err))
		return
	}

	categories, err := h.service.ListByParentID(c.Request.Context(), parentID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, gin.H{"categories": categorydto.ToCategoryListResponse(categories)})
}
