// Package handler implements HTTP request handlers for the example module.
package handler

// 本文件实现 example 模块的 HTTP 处理逻辑。

import (
	"github.com/gin-gonic/gin"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	exampledto "gotribe/internal/api/example/dto"
	exampleservice "gotribe/internal/api/example/service"
	exampleview "gotribe/internal/api/example/view"
	"gotribe/internal/core/middleware"
	"gotribe/internal/request"
	"gotribe/internal/core/response"
)

// Handler 负责处理示例业务单相关的 HTTP 请求。
type Handler struct {
	service *exampleservice.Service
}

// NewHandler 创建示例业务单 HTTP 处理器实例。
func NewHandler(service *exampleservice.Service) *Handler {
	return &Handler{service: service}
}

// Create 处理创建示例业务单请求。
func (h *Handler) Create(c *gin.Context) {
	actor, err := actorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req exampledto.CreateRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	data, err := h.service.Create(c.Request.Context(), middleware.GetProjectID(c), actor, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, toResponse(*data))
}

// List 处理示例业务单列表请求。
func (h *Handler) List(c *gin.Context) {
	actor, err := actorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var query exampledto.ListQuery
	if err := request.BindQuery(c, &query); err != nil {
		response.Error(c, err)
		return
	}

	items, meta, err := h.service.List(c.Request.Context(), middleware.GetProjectID(c), actor, query)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, toListResponse(items), database.Pagination(meta))
}

// Detail 处理示例业务单详情请求。
func (h *Handler) Detail(c *gin.Context) {
	actor, err := actorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	data, err := h.service.Detail(c.Request.Context(), middleware.GetProjectID(c), actor, c.Param("exampleID"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, toResponse(*data))
}

// Update 处理更新示例业务单请求。
func (h *Handler) Update(c *gin.Context) {
	actor, err := actorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	var req exampledto.UpdateRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	data, err := h.service.Update(c.Request.Context(), middleware.GetProjectID(c), actor, c.Param("exampleID"), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, toResponse(*data))
}

// Delete 处理删除示例业务单请求。
func (h *Handler) Delete(c *gin.Context) {
	actor, err := actorFromContext(c)
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), middleware.GetProjectID(c), actor, c.Param("exampleID")); err != nil {
		response.Error(c, err)
		return
	}
	response.NoContent(c)
}

func actorFromContext(c *gin.Context) (exampleview.Actor, error) {
	currentUser, ok := middleware.GetCurrentUser(c)
	if !ok {
		return exampleview.Actor{}, errs.Unauthorized("missing current user context")
	}
	return exampleview.Actor{
		UserID:   currentUser.ID,
		Username: currentUser.Username,
		Nickname: currentUser.Nickname,
	}, nil
}

func toResponse(view exampleview.Example) exampledto.ExampleResponse {
	return exampledto.ExampleResponse{
		ExampleID:   view.ExampleID,
		Name:        view.Name,
		Description: view.Description,
		Status:      view.Status,
		Owner: exampledto.OwnerResponse{
			UserID:   view.Owner.UserID,
			Username: view.Owner.Username,
			Nickname: view.Owner.Nickname,
		},
		PrimaryPost: toPostRefResponse(view.PrimaryPost),
		Posts:       toPostRefResponses(view.Posts),
		CreatedAt:   view.CreatedAt,
		UpdatedAt:   view.UpdatedAt,
	}
}

func toListResponse(items []exampleview.Example) []exampledto.ExampleResponse {
	result := make([]exampledto.ExampleResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toResponse(item))
	}
	return result
}

func toPostRefResponses(items []exampleview.PostRef) []exampledto.PostRefResponse {
	result := make([]exampledto.PostRefResponse, 0, len(items))
	for _, item := range items {
		result = append(result, toPostRefResponse(item))
	}
	return result
}

func toPostRefResponse(item exampleview.PostRef) exampledto.PostRefResponse {
	return exampledto.PostRefResponse{
		PostID: item.PostID,
		Title:  item.Title,
		Type:   item.Type,
		Status: item.Status,
	}
}
