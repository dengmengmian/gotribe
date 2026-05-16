package service

import (
	"context"

	"gotribe/internal/admin/api/dto"
	"gotribe/internal/admin/api/repository"
	"gotribe/internal/core/database"
	"gotribe/internal/model"

	"github.com/casbin/casbin/v2"
)

// Service 接口业务逻辑接口
type Service interface {
	List(ctx context.Context, req *dto.ApiListRequest) ([]*model.Api, int64, error)
	Tree(ctx context.Context) ([]*dto.ApiTreeResponse, error)
	Create(ctx context.Context, actor model.Admin, req *dto.CreateApiRequest) error
	Update(ctx context.Context, actor model.Admin, apiID int64, req *dto.UpdateApiRequest) error
	Delete(ctx context.Context, apiIds []int64) error
}

// service 接口业务逻辑实现
type service struct {
	apiRepo  *repository.Repository
	enforcer *casbin.Enforcer
}

// NewApiService 创建接口服务实例
func NewService(tx *database.TransactionManager, enforcer *casbin.Enforcer) Service {
	return &service{
		apiRepo:  repository.NewRepository(tx, enforcer),
		enforcer: enforcer,
	}
}

// GetApis 获取接口列表
func (s *service) List(ctx context.Context, req *dto.ApiListRequest) ([]*model.Api, int64, error) {
	return s.apiRepo.List(ctx, req)
}

// GetApiTree 获取接口树
func (s *service) Tree(ctx context.Context) ([]*dto.ApiTreeResponse, error) {
	return s.apiRepo.Tree(ctx)
}

// CreateApi 创建接口
func (s *service) Create(ctx context.Context, actor model.Admin, req *dto.CreateApiRequest) error {
	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  actor.Username,
	}
	return s.apiRepo.Create(ctx, &api)
}

// UpdateApiByID 更新接口
func (s *service) Update(ctx context.Context, actor model.Admin, apiID int64, req *dto.UpdateApiRequest) error {
	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  actor.Username,
	}
	return s.apiRepo.Update(ctx, apiID, &api)
}

// BatchDeleteApiByIds 批量删除接口
func (s *service) Delete(ctx context.Context, apiIds []int64) error {
	return s.apiRepo.Delete(ctx, apiIds)
}
