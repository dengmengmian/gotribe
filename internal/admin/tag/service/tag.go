package service

import (
	"context"

	"gotribe/internal/core/database"

	"gotribe/internal/admin/tag/dto"
	"gotribe/internal/admin/tag/repository"
	"gotribe/internal/model"
)

// Service 标签业务逻辑接口
type Service interface {
	Detail(ctx context.Context, id int64) (model.Tag, error)
	List(ctx context.Context, req *dto.TagListRequest) ([]*model.Tag, int64, error)
	Create(ctx context.Context, req *dto.CreateTagRequest) (*model.Tag, error)
	Update(ctx context.Context, id int64, req *dto.CreateTagRequest) error
	Delete(ctx context.Context, ids []int64) error
}

// service 标签业务逻辑实现
type service struct {
	tagRepo *repository.Repository
}

// NewService 创建标签服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		tagRepo: repository.NewRepository(tx),
	}
}

// Detail 根据ID获取标签
func (s *service) Detail(ctx context.Context, id int64) (model.Tag, error) {
	return s.tagRepo.Detail(ctx, id)
}

// List 获取标签列表
func (s *service) List(ctx context.Context, req *dto.TagListRequest) ([]*model.Tag, int64, error) {
	return s.tagRepo.List(ctx, req)
}

// Create 创建标签
func (s *service) Create(ctx context.Context, req *dto.CreateTagRequest) (*model.Tag, error) {
	tag := model.Tag{
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		Color:       req.Color,
	}
	return s.tagRepo.Create(ctx, &tag)
}

// Update 更新标签
func (s *service) Update(ctx context.Context, id int64, req *dto.CreateTagRequest) error {
	oldTag, err := s.tagRepo.Detail(ctx, id)
	if err != nil {
		return err
	}
	oldTag.Title = req.Title
	oldTag.Slug = req.Slug
	oldTag.Description = req.Description
	oldTag.Color = req.Color
	oldTag.Sort = req.Sort
	oldTag.Status = req.Status
	return s.tagRepo.Update(ctx, &oldTag)
}

// Delete 批量删除标签
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.tagRepo.Delete(ctx, ids)
}
