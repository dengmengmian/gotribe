package service

import (
	"context"

	"gotribe/internal/admin/project/dto"
	"gotribe/internal/admin/project/repository"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

// Service 项目业务逻辑接口
type Service interface {
	Detail(ctx context.Context, id int64) (model.Project, error)
	List(ctx context.Context, req *dto.ProjectListRequest) ([]*model.Project, int64, error)
	Create(ctx context.Context, req *dto.CreateProjectRequest) error
	Update(ctx context.Context, id int64, req *dto.CreateProjectRequest) error
	Delete(ctx context.Context, ids []int64) error
}

// service 项目业务逻辑实现
type service struct {
	projectRepo *repository.Repository
}

// NewProjectService 创建项目服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		projectRepo: repository.NewRepository(tx),
	}
}

// Detail 根据ID获取项目
func (s *service) Detail(ctx context.Context, id int64) (model.Project, error) {
	return s.projectRepo.GetProjectByID(ctx, id)
}

// List 获取项目列表
func (s *service) List(ctx context.Context, req *dto.ProjectListRequest) ([]*model.Project, int64, error) {
	return s.projectRepo.List(ctx, req)
}

// Create 创建项目
func (s *service) Create(ctx context.Context, req *dto.CreateProjectRequest) error {
	project := model.Project{
		Name:           req.Name,
		Title:          req.Title,
		Description:    req.Description,
		Keywords:       req.Keywords,
		Domain:         req.Domain,
		PostURL:        req.PostURL,
		ICP:            req.ICP,
		Author:         req.Author,
		Info:           req.Info,
		PublicSecurity: req.PublicSecurity,
		Favicon:        req.Favicon,
		NavImage:       req.NavImage,
		BaiduAnalytics: req.BaiduAnalytics,
		PushToken:      req.PushToken,
	}
	return s.projectRepo.Create(ctx, &project)
}

// Update 更新项目
func (s *service) Update(ctx context.Context, id int64, req *dto.CreateProjectRequest) error {
	oldProject, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return err
	}
	oldProject.Title = req.Title
	oldProject.Description = req.Description
	oldProject.Name = req.Name
	oldProject.Author = req.Author
	oldProject.ICP = req.ICP
	oldProject.Keywords = req.Keywords
	oldProject.Info = req.Info
	oldProject.PostURL = req.PostURL
	oldProject.Domain = req.Domain
	oldProject.PublicSecurity = req.PublicSecurity
	oldProject.Favicon = req.Favicon
	oldProject.NavImage = req.NavImage
	oldProject.BaiduAnalytics = req.BaiduAnalytics
	oldProject.PushToken = req.PushToken
	return s.projectRepo.UpdateProject(ctx, &oldProject)
}

// Delete 批量删除项目
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.projectRepo.Delete(ctx, ids)
}
