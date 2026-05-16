package handler

import (
	"strconv"

	"gotribe/internal/admin/tag/dto"
	"gotribe/internal/admin/tag/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	tagService service.Service
}

// NewHandler 创建标签处理器实例
func NewHandler(tagService service.Service) *Handler {
	return &Handler{tagService: tagService}
}

// Detail 获取标签信息
// @Summary      获取标签信息
// @Description  根据标签ID获取标签详细信息
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        id path int true "标签ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /tag/{id} [get]
// @Security     BearerAuth
func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("标签ID格式错误", nil))
		return
	}
	ctx := c.Request.Context()
	tag, err := h.tagService.Detail(ctx, int64(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	tagInfoDto := dto.ToTagResponse(tag)
	response.OK(c, gin.H{
		"tag": tagInfoDto,
	})
}

// List 获取标签列表
// @Summary      获取标签列表
// @Description  获取所有标签的列表
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        request query dto.TagListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /tag [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.TagListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	tag, total, err := h.tagService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"tags": dto.ToTagListResponse(tag), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Create 创建标签
// @Summary      创建标签
// @Description  创建一个新的标签
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateTagRequest true "创建标签请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /tag [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	tagInfo, err := h.tagService.Create(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, gin.H{"tag": dto.ToTagResponse(*tagInfo)})
}

// Update 更新标签
// @Summary      更新标签
// @Description  根据标签ID更新标签信息
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        id path int true "标签ID"
// @Param        request body dto.CreateTagRequest true "标签信息"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /tag/{id} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("标签ID格式错误", nil))
		return
	}
	ctx := c.Request.Context()
	err = h.tagService.Update(ctx, int64(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除标签
// @Summary      批量删除标签
// @Description  根据标签ID列表批量删除标签
// @Tags         标签管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteTagsRequest true "标签ID列表"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /tag [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteTagsRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.tagService.Delete(ctx, req.Ids)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
