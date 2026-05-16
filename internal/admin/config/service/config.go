package service

import (
	"context"

	"gotribe/internal/admin/config/dto"
	"gotribe/internal/admin/config/repository"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

// Service 配置业务逻辑接口
type Service interface {
	Detail(ctx context.Context, id int64) (model.Config, error)
	DetailByAlias(ctx context.Context, alias string) (model.Config, error)
	List(ctx context.Context, req *dto.ConfigListRequest) ([]*model.Config, int64, error)
	Create(ctx context.Context, req *dto.CreateConfigRequest) error
	Update(ctx context.Context, id int64, req *dto.UpdateConfigRequest) error
	Delete(ctx context.Context, ids []int64) error
}

// service 配置业务逻辑实现
type service struct {
	configRepo *repository.Repository
}

// NewConfigService 创建配置服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		configRepo: repository.NewRepository(tx),
	}
}

// Detail 根据自增 ID 获取配置
func (s *service) Detail(ctx context.Context, id int64) (model.Config, error) {
	return s.configRepo.GetConfigByID(ctx, id)
}

// DetailByAlias 根据别名获取配置
func (s *service) DetailByAlias(ctx context.Context, alias string) (model.Config, error) {
	return s.configRepo.GetConfigByAlias(ctx, alias)
}

// List 获取配置列表
func (s *service) List(ctx context.Context, req *dto.ConfigListRequest) ([]*model.Config, int64, error) {
	return s.configRepo.List(ctx, req)
}

// Create 创建配置
func (s *service) Create(ctx context.Context, req *dto.CreateConfigRequest) error {
	config := model.Config{
		ProjectID:   req.ProjectID,
		Alias:       req.Alias,
		Title:       req.Title,
		MDContent:   req.MDContent,
		Description: req.Description,
		Type:        req.Type,
		Info:        req.Info,
	}
	return s.configRepo.Create(ctx, &config)
}

// Update 更新配置
func (s *service) Update(ctx context.Context, id int64, req *dto.UpdateConfigRequest) error {
	oldConfig, err := s.configRepo.GetConfigByID(ctx, id)
	if err != nil {
		return err
	}
	oldConfig.Title = req.Title
	oldConfig.Description = req.Description
	oldConfig.Info = req.Info
	oldConfig.ProjectID = req.ProjectID
	oldConfig.MDContent = req.MDContent
	return s.configRepo.UpdateConfig(ctx, &oldConfig)
}

// Delete 批量删除配置
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.configRepo.Delete(ctx, ids)
}
