package handler

import (
	"strconv"

	"gotribe/internal/admin/api/dto"
	"gotribe/internal/admin/api/service"
	"gotribe/internal/admin/common"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	apiService service.Service
}

// NewHandler 创建接口处理器实例
func NewHandler(apiService service.Service) *Handler {
	return &Handler{apiService: apiService}
}

// List 获取接口列表
// @Summary      获取接口列表
// @Description  获取所有接口的列表
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        request query dto.ApiListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /api/list [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.ApiListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	apis, total, err := h.apiService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{
		"apis": dto.ToApiListResponse(apis), "total": total,
	}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Tree 获取接口树
// @Summary      获取接口树
// @Description  获取树形结构的接口列表(按接口Category字段分类)
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /api/tree [get]
// @Security     BearerAuth
func (h *Handler) Tree(c *gin.Context) {
	ctx := c.Request.Context()
	tree, err := h.apiService.Tree(ctx)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{
		"apis": tree,
	})
}

// Create 创建接口
// @Summary      创建接口
// @Description  创建一个新的接口
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateApiRequest true "创建接口请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /api/create [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateApiRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.apiService.Create(ctx, actor, &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, nil)
}

// Update 更新接口
// @Summary      更新接口
// @Description  根据接口ID更新接口信息
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        apiID path string true "接口ID"
// @Param        request body dto.UpdateApiRequest true "更新接口请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /api/update/{apiID} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdateApiRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	apiID, _ := strconv.Atoi(c.Param("apiID"))
	if apiID <= 0 {
		response.Error(c, errs.BadRequest("接口ID不正确", nil))
		return
	}

	ctx := c.Request.Context()
	actor, err := common.CurrentAdmin(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	err = h.apiService.Update(ctx, actor, int64(apiID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, nil)
}

// Delete 批量删除接口
// @Summary      批量删除接口
// @Description  根据接口ID列表批量删除接口
// @Tags         接口管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteApiRequest true "删除接口请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /api/delete/batch [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteApiRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.apiService.Delete(ctx, req.ApiIds)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
