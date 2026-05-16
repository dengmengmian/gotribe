// Package handler implements HTTP request handlers for the post module.
package handler

// 本文件实现帖子模块的 HTTP 处理逻辑。

import (
	"github.com/gin-gonic/gin"
	"gotribe/internal/core/middleware"
	postdto "gotribe/internal/api/post/dto"
	postservice "gotribe/internal/api/post/service"
	"gotribe/internal/request"
	"gotribe/internal/core/response"
)

// Handler 负责处理文章相关的 HTTP 请求。
type Handler struct {
	service *postservice.Service
}

// NewHandler 创建文章 HTTP 处理器实例。
func NewHandler(service *postservice.Service) *Handler {
	return &Handler{service: service}
}

// List 处理文章列表查询请求。
func (h *Handler) List(c *gin.Context) {
	var query postdto.ListQuery
	if err := request.BindQuery(c, &query); err != nil {
		response.Error(c, err)
		return
	}

	items, meta, err := h.service.List(c.Request.Context(), middleware.GetProjectID(c), query)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, items, meta)
}

// Detail 处理文章详情查询请求。
func (h *Handler) Detail(c *gin.Context) {
	var query postdto.DetailQuery
	if err := request.BindQuery(c, &query); err != nil {
		response.Error(c, err)
		return
	}

	data, err := h.service.Detail(c.Request.Context(), middleware.GetProjectID(c), c.Param("postID"), query.Password)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, data)
}
