package handler

import (
	"strconv"

	"gotribe/internal/admin/post/dto"
	"gotribe/internal/admin/post/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	postService service.Service
}

// NewHandler 创建内容处理器实例
func NewHandler(postService service.Service) *Handler {
	return &Handler{postService: postService}
}

// Detail 获取内容信息
// @Summary      获取内容信息
// @Description  根据内容ID获取内容详细信息
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        id path int64 true "内容ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /post/{id} [get]
// @Security     BearerAuth
func (h *Handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的 id", nil))
		return
	}
	post, err := h.postService.Detail(ctx, int64(postID))
	if err != nil {
		response.Error(c, err)
		return
	}
	postResponse := dto.ToPostResponse(&post)
	response.OK(c, gin.H{
		"post": postResponse,
	})
}

// List 获取内容列表
// @Summary      获取内容列表
// @Description  获取所有内容的列表
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        request query dto.PostListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /post [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.PostListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	post, total, err := h.postService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"posts": dto.ToPostListResponse(post), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Create 创建内容
// @Summary      创建内容
// @Description  创建一个新的内容
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        request body dto.CreatePostRequest true "创建内容请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /post [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.postService.Create(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// Update 更新内容
// @Summary      更新内容
// @Description  根据内容ID更新内容信息
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        id path int64 true "内容ID"
// @Param        request body dto.UpdatePostRequest true "更新内容请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /post/{id} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdatePostRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的 id", nil))
		return
	}
	err = h.postService.Update(ctx, int64(postID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除内容
// @Summary      批量删除内容
// @Description  根据内容ID列表批量删除内容
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeletePostsRequest true "删除内容请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /post [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeletePostsRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.postService.Delete(ctx, req.PostIds)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Publish 发布内容
// @Summary      发布内容
// @Description  根据内容ID发布内容
// @Tags         内容管理
// @Accept       json
// @Produce      json
// @Param        id path int64 true "内容ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /post/{id} [put]
// @Security     BearerAuth
func (h *Handler) Publish(c *gin.Context) {
	ctx := c.Request.Context()
	postID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的 id", nil))
		return
	}
	err = h.postService.Publish(ctx, int64(postID))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
