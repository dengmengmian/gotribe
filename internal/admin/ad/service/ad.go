package service

import (
	"context"

	"gotribe/internal/admin/ad/dto"
	"gotribe/internal/admin/ad/repository"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

// Service 广告业务逻辑接口
type Service interface {
	Detail(ctx context.Context, id int64) (model.Ad, error)
	List(ctx context.Context, req *dto.AdListRequest) ([]*model.Ad, int64, error)
	Create(ctx context.Context, req *dto.CreateAdRequest) error
	Update(ctx context.Context, id int64, req *dto.UpdateAdRequest) error
	Delete(ctx context.Context, ids []int64) error
}

// service 广告业务逻辑实现
type service struct {
	adRepo *repository.Repository
}

// NewAdService 创建广告服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		adRepo: repository.NewRepository(tx),
	}
}

// Detail 根据ID获取广告
func (s *service) Detail(ctx context.Context, id int64) (model.Ad, error) {
	return s.adRepo.Detail(ctx, id)
}

// List 获取广告列表
func (s *service) List(ctx context.Context, req *dto.AdListRequest) ([]*model.Ad, int64, error) {
	return s.adRepo.List(ctx, req)
}

// Create 创建广告
func (s *service) Create(ctx context.Context, req *dto.CreateAdRequest) error {
	ad := model.Ad{
		SceneID:     req.SceneID,
		URL:         req.URL,
		URLType:     req.URLType,
		Image:       req.Image,
		Sort:        req.Sort,
		Status:      req.Status,
		Title:       req.Title,
		Video:       req.Video,
		Ext:         req.Ext,
		Description: req.Description,
	}
	return s.adRepo.Create(ctx, &ad)
}

// Update 更新广告
func (s *service) Update(ctx context.Context, id int64, req *dto.UpdateAdRequest) error {
	oldAd, err := s.adRepo.Detail(ctx, id)
	if err != nil {
		return err
	}
	oldAd.Title = req.Title
	oldAd.Description = req.Description
	oldAd.Ext = req.Ext
	oldAd.Image = req.Image
	oldAd.Sort = req.Sort
	oldAd.URL = req.URL
	oldAd.URLType = req.URLType
	oldAd.Status = req.Status
	oldAd.Video = req.Video
	oldAd.SceneID = req.SceneID
	return s.adRepo.Update(ctx, &oldAd)
}

// Delete 批量删除广告
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.adRepo.Delete(ctx, ids)
}
