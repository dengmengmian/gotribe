package service

import (
	"context"

	"gotribe/internal/admin/point/dto"
	pointRepository "gotribe/internal/admin/point/repository"
	"gotribe/internal/core/database"
	"gotribe/internal/model"
)

// Service 积分业务逻辑接口。
type Service interface {
	List(ctx context.Context, req *dto.PointLogListRequest) ([]*model.PointLog, int64, error)
	Create(ctx context.Context, req *dto.CreatePointLogRequest) error
	Detail(ctx context.Context, id int64) (*model.PointLog, error)
	Update(ctx context.Context, id int64, req *dto.UpdatePointLogRequest) error
	Delete(ctx context.Context, ids []int64) error
}

type service struct {
	pointLogRepo *pointRepository.Repository
	tx           *database.TransactionManager
}

// NewService 创建积分服务实例。
func NewService(tx *database.TransactionManager) Service {
	return &service{
		pointLogRepo: pointRepository.NewRepository(tx),
		tx:           tx,
	}
}

// List 获取积分列表。
func (s *service) List(ctx context.Context, req *dto.PointLogListRequest) ([]*model.PointLog, int64, error) {
	return s.pointLogRepo.List(ctx, req)
}

// Create 创建积分。
func (s *service) Create(ctx context.Context, req *dto.CreatePointLogRequest) error {
	return s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.pointLogRepo.Create(txCtx, req)
	})
}

// Detail 获取单条积分记录。
func (s *service) Detail(ctx context.Context, id int64) (*model.PointLog, error) {
	return s.pointLogRepo.Detail(ctx, id)
}

// Update 更新积分记录。
func (s *service) Update(ctx context.Context, id int64, req *dto.UpdatePointLogRequest) error {
	log := &model.PointLog{}
	if req.Point != 0 {
		log.Points = int64(req.Point * 100)
	}
	if req.UserID != 0 {
		log.UserID = req.UserID
	}
	if req.ProjectID != 0 {
		log.ProjectID = req.ProjectID
	}
	return s.pointLogRepo.Update(ctx, id, log)
}

// Delete 批量删除积分记录。
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.pointLogRepo.Delete(ctx, ids)
}
