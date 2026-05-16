package service

import (
	"context"

	"gotribe/internal/admin/column/dto"
	"gotribe/internal/admin/column/repository"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

// Service 专栏业务逻辑接口
type Service interface {
	Detail(ctx context.Context, id int64) (model.Column, error)
	List(ctx context.Context, req *dto.ColumnListRequest) ([]*model.Column, int64, error)
	Create(ctx context.Context, req *dto.CreateColumnRequest) error
	Update(ctx context.Context, id int64, req *dto.UpdateColumnRequest) error
	Delete(ctx context.Context, ids []int64) error
}

// service 专栏业务逻辑实现
type service struct {
	columnRepo *repository.Repository
}

// NewService 创建专栏服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		columnRepo: repository.NewRepository(tx),
	}
}

// Detail 根据ID获取专栏
func (s *service) Detail(ctx context.Context, id int64) (model.Column, error) {
	return s.columnRepo.Detail(ctx, id)
}

// List 获取专栏列表
func (s *service) List(ctx context.Context, req *dto.ColumnListRequest) ([]*model.Column, int64, error) {
	return s.columnRepo.List(ctx, req)
}

// Create 创建专栏
func (s *service) Create(ctx context.Context, req *dto.CreateColumnRequest) error {
	column := model.Column{
		Title:       req.Title,
		Description: req.Description,
		Info:        req.Info,
		Icon:        req.Icon,
		ProjectID:   req.ProjectID,
	}
	return s.columnRepo.Create(ctx, &column)
}

// Update 更新专栏
func (s *service) Update(ctx context.Context, id int64, req *dto.UpdateColumnRequest) error {
	oldColumn, err := s.columnRepo.Detail(ctx, id)
	if err != nil {
		return err
	}
	oldColumn.Title = req.Title
	oldColumn.Description = req.Description
	oldColumn.Info = req.Info
	oldColumn.Icon = req.Icon
	oldColumn.ProjectID = req.ProjectID
	return s.columnRepo.Update(ctx, &oldColumn)
}

// Delete 批量删除专栏
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.columnRepo.Delete(ctx, ids)
}
