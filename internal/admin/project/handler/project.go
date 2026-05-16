package handler

import (
	"gotribe/internal/admin/project/dto"
	"gotribe/internal/admin/project/service"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"
	"gotribe/internal/request"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	projectService service.Service
}

// NewHandler 创建项目处理器实例
func NewHandler(projectService service.Service) *Handler {
	return &Handler{projectService: projectService}
}

// Detail 获取项目信息
// @Summary      获取项目信息
// @Description  根据项目ID获取项目详细信息
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        id path int64 true "项目ID"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /project/{id} [get]
// @Security     BearerAuth
func (h *Handler) Detail(c *gin.Context) {
	ctx := c.Request.Context()
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的 id", nil))
		return
	}
	project, err := h.projectService.Detail(ctx, int64(projectID))
	if err != nil {
		response.Error(c, err)
		return
	}
	projectInfoDto := dto.ToProjectResponse(&project)
	response.OK(c, gin.H{
		"project": projectInfoDto,
	})
}

// List 获取项目列表
// @Summary      获取项目列表
// @Description  获取所有项目的列表，支持分页和筛选
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        request query dto.ProjectListRequest false "查询参数"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /project [get]
// @Security     BearerAuth
func (h *Handler) List(c *gin.Context) {
	var req dto.ProjectListRequest
	if err := request.BindQuery(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	project, total, err := h.projectService.List(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, 200, gin.H{"projects": dto.ToProjectListResponse(project), "total": total}, gin.H{"page": req.PageNum, "per_page": req.PageSize, "total": total})
}

// Create 创建项目
// @Summary      创建项目
// @Description  创建一个新的项目
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateProjectRequest true "创建项目请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /project [post]
// @Security     BearerAuth
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.projectService.Create(ctx, &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, nil)
}

// Update 更新项目
// @Summary      更新项目
// @Description  根据项目ID更新项目信息
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        id path int64 true "项目ID"
// @Param        request body dto.CreateProjectRequest true "更新项目请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /project/{id} [patch]
// @Security     BearerAuth
func (h *Handler) Update(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	projectID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errs.BadRequest("无效的 id", nil))
		return
	}

	err = h.projectService.Update(ctx, int64(projectID), &req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

// Delete 批量删除项目
// @Summary      批量删除项目
// @Description  根据项目ID列表批量删除项目
// @Tags         项目管理
// @Accept       json
// @Produce      json
// @Param        request body dto.DeleteProjectsRequest true "删除项目请求"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /project [delete]
// @Security     BearerAuth
func (h *Handler) Delete(c *gin.Context) {
	var req dto.DeleteProjectsRequest
	if err := request.BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}

	ctx := c.Request.Context()
	err := h.projectService.Delete(ctx, req.ProjectIDs)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
