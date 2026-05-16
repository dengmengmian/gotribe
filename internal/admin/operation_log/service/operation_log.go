package service

import (
	"context"

	"gotribe/internal/model"
	"gotribe/internal/admin/operation_log/dto"
	"gotribe/internal/admin/operation_log/repository"

	"go.uber.org/zap"
	"gotribe/internal/core/database"
)

// Service 操作日志业务逻辑接口
type Service interface {
	List(ctx context.Context, req *dto.OperationLogListRequest) ([]model.OperationLog, int64, error)
	Delete(ctx context.Context, ids []int64) error
}

// service 操作日志业务逻辑实现
type service struct {
	operationLogRepo *repository.Repository
}

// NewService 创建操作日志服务实例
func NewService(tx *database.TransactionManager, log *zap.SugaredLogger) Service {
	return &service{
		operationLogRepo: repository.NewRepository(tx, log),
	}
}

// List 获取操作日志列表
func (s *service) List(ctx context.Context, req *dto.OperationLogListRequest) ([]model.OperationLog, int64, error) {
	return s.operationLogRepo.List(ctx, req)
}

// Delete 批量删除操作日志
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.operationLogRepo.Delete(ctx, ids)
}
