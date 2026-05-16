package handler

import (
	"strconv"

	"gotribe/internal/admin/job/dto"
	"gotribe/internal/admin/jobs"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/response"

	"github.com/gin-gonic/gin"
)

// Handler 任务处理器
type Handler struct{}

// NewHandler 创建任务处理器实例
func NewHandler() *Handler {
	return &Handler{}
}

// ListJobs 列出所有任务
// @Summary      获取任务列表
// @Description  获取系统中所有定时任务的列表
// @Tags         任务管理
// @Accept       json
// @Produce      json
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /job [get]
// @Security     BearerAuth
func (h *Handler) ListJobs(ctx *gin.Context) {
	registry := jobs.GetGlobalRegistry()
	jobList := registry.ListJobs()

	var jobsVO []dto.JobVO
	for _, job := range jobList {
		jobVO := dto.JobVO{
			Name:        job.Name(),
			Description: job.Description(),
			Schedule:    job.Schedule(),
			Enabled:     job.IsEnabled(),
			Timeout:     job.Timeout().String(),
			RetryCount:  job.RetryCount(),
		}
		jobsVO = append(jobsVO, jobVO)
	}

	response.OK(ctx, gin.H{"jobs": jobsVO})
}

// GetJobStatus 获取任务状态
// @Summary      获取任务状态
// @Description  根据任务名称获取指定任务的运行状态
// @Tags         任务管理
// @Accept       json
// @Produce      json
// @Param        name path string true "任务名称"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /job/{name}/status [get]
// @Security     BearerAuth
func (h *Handler) GetJobStatus(ctx *gin.Context) {
	jobName := ctx.Param("name")
	if jobName == "" {
		response.Error(ctx, errs.BadRequest("任务名称不能为空", nil))
		return
	}

	registry := jobs.GetGlobalRegistry()
	status, err := registry.GetJobStatus(jobName)
	if err != nil {
		response.Error(ctx, errs.Internal("任务状态获取失败", err))
		return
	}

	response.OK(ctx, gin.H{"status": status})
}

// GetJobHistory 获取任务执行历史
// @Summary      获取任务执行历史
// @Description  根据任务名称获取指定任务的执行历史记录
// @Tags         任务管理
// @Accept       json
// @Produce      json
// @Param        name path string true "任务名称"
// @Param        limit query int false "限制返回记录数" default(10)
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /job/{name}/history [get]
// @Security     BearerAuth
func (h *Handler) GetJobHistory(ctx *gin.Context) {
	jobName := ctx.Param("name")
	if jobName == "" {
		response.Error(ctx, errs.BadRequest("任务名称不能为空", nil))
		return
	}

	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	registry := jobs.GetGlobalRegistry()
	history, err := registry.GetJobHistory(jobName, limit)
	if err != nil {
		response.Error(ctx, errs.Internal("任务历史获取失败", err))
		return
	}

	response.OK(ctx, gin.H{"history": history})
}

// EnableJob 启用任务
// @Summary      启用任务
// @Description  根据任务名称启用指定的定时任务
// @Tags         任务管理
// @Accept       json
// @Produce      json
// @Param        name path string true "任务名称"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /job/{name}/enable [post]
// @Security     BearerAuth
func (h *Handler) EnableJob(ctx *gin.Context) {
	jobName := ctx.Param("name")
	if jobName == "" {
		response.Error(ctx, errs.BadRequest("任务名称不能为空", nil))
		return
	}

	registry := jobs.GetGlobalRegistry()
	if err := registry.EnableJob(jobName); err != nil {
		response.Error(ctx, errs.Internal("启用任务失败", err))
		return
	}

	response.OK(ctx, nil)
}

// DisableJob 禁用任务
// @Summary      禁用任务
// @Description  根据任务名称禁用指定的定时任务
// @Tags         任务管理
// @Accept       json
// @Produce      json
// @Param        name path string true "任务名称"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Router       /job/{name}/disable [post]
// @Security     BearerAuth
func (h *Handler) DisableJob(ctx *gin.Context) {
	jobName := ctx.Param("name")
	if jobName == "" {
		response.Error(ctx, errs.BadRequest("任务名称不能为空", nil))
		return
	}

	registry := jobs.GetGlobalRegistry()
	if err := registry.DisableJob(jobName); err != nil {
		response.Error(ctx, errs.Internal("禁用任务失败", err))
		return
	}

	response.OK(ctx, nil)
}
