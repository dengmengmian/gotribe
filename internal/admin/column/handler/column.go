package handler

import (
	"strconv"

	"gotribe/internal/admin/column/dto"
	"gotribe/internal/admin/column/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	columnService service.Service
}

// NewHandler 创建专栏处理器实例
func NewHandler(columnService service.Service) *Handler {
	return &Handler{columnService: columnService}
}

// Detail 获取专栏信息
// @Summary      获取专栏详情
// @Description  根据专栏ID获取专栏详细信息
// @Tags         专栏管理
// @Accept       json
// @Produce      json
// @Param        id path int true "专栏ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /column/{id} [get]
// @Security     BearerAuth
func (h *Handler) Detail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("专栏ID格式错误", err))
		return
	}
	ctx := c.Request.Context()
	column, err := h.columnService.Detail(ctx, int64(id))
	if err != nil {
		response.Error(c, err)
		return
	}
	columnInfoDto := dto.ToColumnResponse(column)
	response.OK(c, gin.H{
		"column": columnInfoDto,
	})
}

// List 获取专栏列表
// @Summary      获取专栏列表
// @Description  根据查询条件获取专栏列表，支持分页
// @Tags         专栏管理
// @Accept       json
// @Produce      json
// @Param        request query dto.ColumnListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /column [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.ColumnListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	// 获取
	column, total, err := h.columnService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"columns": dto.ToColumnListResponse(column), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Create 创建专栏
// @Summary      创建专栏
// @Description  创建新的专栏
// @Tags         专栏管理
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateColumnRequest true "专栏信息"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /column [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateColumnRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.columnService.Create(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// Update 更新专栏
// @Summary      更新专栏
// @Description  根据专栏ID更新专栏信息
// @Tags         专栏管理
// @Accept       json
// @Produce      json
// @Param        id path int true "专栏ID"
// @Param        request body dto.UpdateColumnRequest true "专栏信息"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /column/{id} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("专栏ID格式错误", err))
		return
	}
	var req dto.UpdateColumnRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	// 更新专栏
	err = h.columnService.Update(ctx, int64(id), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除专栏
// @Summary      批量删除专栏
// @Description  根据专栏ID列表批量删除专栏
// @Tags         专栏管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteColumnsRequest true "专栏ID列表"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /column [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteColumnsRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.columnService.Delete(ctx, req.Ids)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
