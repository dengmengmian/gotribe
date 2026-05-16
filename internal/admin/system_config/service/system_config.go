package service

import (
	"context"

	"gotribe/internal/model"
	"gotribe/internal/admin/system_config/dto"
	"gotribe/internal/admin/system_config/repository"

	"gotribe/internal/core/database"
)

// Service 系统配置业务逻辑接口
type Service interface {
	Detail(ctx context.Context) (model.SystemConfig, error)
	Update(ctx context.Context, req *dto.CreateSystemConfigRequest) error
}

// service 系统配置业务逻辑实现
type service struct {
	systemConfigRepo *repository.Repository
}

// NewSystemConfigService 创建系统配置服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		systemConfigRepo: repository.NewRepository(tx),
	}
}

// Detail 获取系统配置
func (s *service) Detail(ctx context.Context) (model.SystemConfig, error) {
	return s.systemConfigRepo.Detail(ctx)
}

// Update 更新系统配置
func (s *service) Update(ctx context.Context, req *dto.CreateSystemConfigRequest) error {
	oldSystemConfig, err := s.systemConfigRepo.Detail(ctx)
	if err != nil {
		return err
	}
	oldSystemConfig.Title = req.Title
	oldSystemConfig.Content = req.Content
	oldSystemConfig.Footer = req.Footer
	oldSystemConfig.Icon = req.Icon
	oldSystemConfig.Logo = req.Logo
	return s.systemConfigRepo.Update(ctx, &oldSystemConfig)
}
