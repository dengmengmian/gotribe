package handler

import (
	"strconv"

	"gotribe/internal/admin/ad/dto"
	"gotribe/internal/admin/ad/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	adService service.Service
}

// NewHandler 创建广告处理器实例
func NewHandler(adService service.Service) *Handler {
	return &Handler{adService: adService}
}

// GetAdInfo 获取当前广告信息
// @Summary 获取广告详情
// @Description 根据广告ID获取广告详细信息
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param id path int true "广告ID"
// @Success 200 {object} response.Response{data=map[string]dto.AdResponse} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/{id} [get]
// @Security BearerAuth
func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的ID", err))
		return
	}
	ad, err := h.adService.Detail(c.Request.Context(), int64(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	adInfoDto := dto.ToAdInfoResponse(ad)
	response.OK(c, gin.H{
		"ad": adInfoDto,
	})
}

// GetAds 获取广告列表
// @Summary 获取广告列表
// @Description 根据查询条件获取广告列表，支持分页
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param sceneID query string false "场景ID"
// @Param title query string false "广告标题"
// @Param status query int false "状态(1:启用 2:禁用)"
// @Param pageNum query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=map[string]interface{}} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad [get]
// @Security BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.AdListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	// 获取
	ad, total, err := h.adService.List(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"ads": dto.ToAdListResponse(ad), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// CreateAd 创建广告
// @Summary 创建广告
// @Description 创建新的广告
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param ad body dto.CreateAdRequest true "广告信息"
// @Success 200 {object} response.Response "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad [post]
// @Security BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateAdRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	err := h.adService.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// UpdateAdByID 更新广告
// @Summary 更新广告
// @Description 根据广告ID更新广告信息
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param id path int true "广告ID"
// @Param ad body dto.UpdateAdRequest true "广告信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad/{id} [patch]
// @Security BearerAuth
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的ID", err))
		return
	}
	var req dto.UpdateAdRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	err = h.adService.Update(c.Request.Context(), int64(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// BatchDeleteAdByIds 批量删除广告
// @Summary 批量删除广告
// @Description 根据广告ID列表批量删除广告
// @Tags 广告管理
// @Accept json
// @Produce json
// @Param ads body dto.DeleteAdsRequest true "广告ID列表"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /ad [delete]
// @Security BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteAdsRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	// 前端传来的广告ID
	err := h.adService.Delete(c.Request.Context(), req.Ids)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
