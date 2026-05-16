package service

import (
	"context"

	"gotribe/internal/admin/ad_scene/dto"
	"gotribe/internal/admin/ad_scene/repository"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

// Service 推广场景业务逻辑接口
type Service interface {
	Detail(ctx context.Context, id int64) (model.AdScene, error)
	List(ctx context.Context, req *dto.AdSceneListRequest) ([]*model.AdScene, int64, error)
	Create(ctx context.Context, req *dto.CreateAdSceneRequest) error
	Update(ctx context.Context, id int64, req *dto.UpdateAdSceneRequest) error
	Delete(ctx context.Context, ids []int64) error
}

// service 推广场景业务逻辑实现
type service struct {
	adSceneRepo *repository.Repository
}

// NewAdSceneService 创建推广场景服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		adSceneRepo: repository.NewRepository(tx),
	}
}

// Detail 根据ID获取推广场景
func (s *service) Detail(ctx context.Context, id int64) (model.AdScene, error) {
	return s.adSceneRepo.Detail(ctx, id)
}

// List 获取推广场景列表
func (s *service) List(ctx context.Context, req *dto.AdSceneListRequest) ([]*model.AdScene, int64, error) {
	return s.adSceneRepo.List(ctx, req)
}

// Create 创建推广场景
func (s *service) Create(ctx context.Context, req *dto.CreateAdSceneRequest) error {
	adScene := model.AdScene{
		ProjectID:   req.ProjectID,
		Title:       req.Title,
		Description: req.Description,
	}
	return s.adSceneRepo.Create(ctx, &adScene)
}

// Update 更新推广场景
func (s *service) Update(ctx context.Context, id int64, req *dto.UpdateAdSceneRequest) error {
	oldAdScene, err := s.adSceneRepo.Detail(ctx, id)
	if err != nil {
		return err
	}
	oldAdScene.Title = req.Title
	oldAdScene.Description = req.Description
	oldAdScene.ProjectID = req.ProjectID
	return s.adSceneRepo.Update(ctx, &oldAdScene)
}

// Delete 批量删除推广场景
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.adSceneRepo.Delete(ctx, ids)
}
