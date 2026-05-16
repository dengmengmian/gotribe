package handler

import (
	"strconv"

	"gotribe/internal/admin/config/dto"
	"gotribe/internal/admin/config/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	configService service.Service
}

// NewHandler 创建配置处理器实例
func NewHandler(configService service.Service) *Handler {
	return &Handler{configService: configService}
}

// Detail 获取配置信息
// @Summary      获取配置信息
// @Description  根据配置ID获取配置详细信息
// @Tags         配置管理
// @Accept       json
// @Produce      json
// @Param        configID path string true "配置ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /config/{configID} [get]
// @Security     BearerAuth
func (h *Handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("配置ID格式错误", err))
		return
	}
	config, err := h.configService.Detail(ctx, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	configInfoDto := dto.ToConfigResponse(config)
	response.OK(c, gin.H{
		"config": configInfoDto,
	})
}

// DetailByAlias 根据别名获取配置信息
func (h *Handler) DetailByAlias(c *gin.Context) {
	ctx := c.Request.Context()
	config, err := h.configService.DetailByAlias(ctx, c.Param("alias"))
	if err != nil {
		response.Error(c, err)
		return
	}
	configInfoDto := dto.ToConfigResponse(config)
	response.OK(c, gin.H{
		"config": configInfoDto,
	})
}

// List 获取配置列表
// @Summary      获取配置列表
// @Description  获取所有配置的列表
// @Tags         配置管理
// @Accept       json
// @Produce      json
// @Param        request query dto.ConfigListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /config [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.ConfigListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	// 获取
	config, total, err := h.configService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"configs": dto.ToConfigListResponse(config), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Create 创建配置
// @Summary      创建配置
// @Description  创建一个新的配置
// @Tags         配置管理
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateConfigRequest true "创建配置请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /config [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateConfigRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.configService.Create(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// Update 更新配置
// @Summary      更新配置
// @Description  根据配置ID更新配置信息
// @Tags         配置管理
// @Accept       json
// @Produce      json
// @Param        configID path string true "配置ID"
// @Param        request body dto.UpdateConfigRequest true "配置信息"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /config/{configID} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdateConfigRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("配置ID格式错误", err))
		return
	}
	// 更新配置
	err = h.configService.Update(ctx, id, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除配置
// @Summary      批量删除配置
// @Description  根据配置ID列表批量删除配置
// @Tags         配置管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteConfigsRequest true "配置ID列表"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /config [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteConfigsRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.configService.Delete(ctx, req.Ids)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
