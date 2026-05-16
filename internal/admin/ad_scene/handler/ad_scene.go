package handler

import (
	"strconv"

	"gotribe/internal/admin/ad_scene/dto"
	"gotribe/internal/admin/ad_scene/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	adSceneService service.Service
}

// NewHandler 创建推广场景处理器实例
func NewHandler(adSceneService service.Service) *Handler {
	return &Handler{adSceneService: adSceneService}
}

// Detail 获取当前推广场景信息
// @Summary 获取推广场景详情
// @Description 根据推广场景ID获取推广场景详细信息
// @Tags 推广场景管理
// @Accept json
// @Produce json
// @Param id path int true "推广场景ID"
// @Success 200 {object} response.Response{data=map[string]dto.AdSceneResponse} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/scene/{id} [get]
// @Security BearerAuth
func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的ID", err))
		return
	}
	adScene, err := h.adSceneService.Detail(c.Request.Context(), int64(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	adSceneInfoDto := dto.ToAdSceneInfoResponse(adScene)
	response.OK(c, gin.H{
		"adScene": adSceneInfoDto,
	})
}

// List 获取推广场景列表
// @Summary 获取推广场景列表
// @Description 根据查询条件获取推广场景列表，支持分页
// @Tags 推广场景管理
// @Accept json
// @Produce json
// @Param ProjectID query int64 false "项目ID"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/scene [get]
// @Security BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.AdSceneListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	// 获取
	adScene, total, err := h.adSceneService.List(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"ad_scenes": dto.ToAdSceneListResponse(adScene), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Create 创建推广场景
// @Summary 创建推广场景
// @Description 创建新的推广场景
// @Tags 推广场景管理
// @Accept json
// @Produce json
// @Param adScene body dto.CreateAdSceneRequest true "推广场景信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/scene [post]
// @Security BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateAdSceneRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	err := h.adSceneService.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// Update 更新推广场景
// @Summary 更新推广场景
// @Description 根据推广场景ID更新推广场景信息
// @Tags 推广场景管理
// @Accept json
// @Produce json
// @Param id path int true "推广场景ID"
// @Param adScene body dto.UpdateAdSceneRequest true "推广场景信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/scene/{id} [patch]
// @Security BearerAuth
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的ID", err))
		return
	}
	var req dto.UpdateAdSceneRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	err = h.adSceneService.Update(c.Request.Context(), int64(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除推广场景
// @Summary 批量删除推广场景
// @Description 根据推广场景ID列表批量删除推广场景
// @Tags 推广场景管理
// @Accept json
// @Produce json
// @Param ad_scenes body dto.DeleteAdScenesRequest true "推广场景ID列表"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/scene [delete]
// @Security BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteAdScenesRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	// 前端传来的推广场景ID
	err := h.adSceneService.Delete(c.Request.Context(), req.Ids)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
