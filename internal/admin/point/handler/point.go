package handler

import (
	"gotribe/internal/admin/point/dto"
	"gotribe/internal/admin/point/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	pointService service.Service
}

// NewHandler 创建积分处理器实例
func NewHandler(pointService service.Service) *Handler {
	return &Handler{pointService: pointService}
}

// GetPoints 获取积分列表
// @Summary      获取积分列表
// @Description  根据查询条件获取积分列表，支持分页
// @Tags         积分管理
// @Accept       json
// @Produce      json
// @Param        userID query string false "用户ID"
// @Param        projectID query string false "项目ID"
// @Param        pageNum query int false "页码"
// @Param        pageSize query int false "每页数量"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /point [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.PointLogListRequest
	// 参数绑定
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	// 获取
	point, total, err := h.pointService.List(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"points": dto.ToPointListResponse(point), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// CreatePoint 创建积分
// @Summary      创建积分
// @Description  为用户创建积分记录
// @Tags         积分管理
// @Accept       json
// @Produce      json
// @Param        request body dto.CreatePointLogRequest true "积分信息"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /point [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreatePointLogRequest
	// 参数绑定
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	err := h.pointService.Create(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// Detail 获取积分详情。
func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的ID", err))
		return
	}
	log, err := h.pointService.Detail(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, gin.H{"point": log})
}

// Update 更新积分记录。
func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdatePointLogRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的ID", err))
		return
	}
	if err := h.pointService.Update(c.Request.Context(), id, &req); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除积分记录。
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeletePointLogRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.pointService.Delete(c.Request.Context(), req.Ids); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
