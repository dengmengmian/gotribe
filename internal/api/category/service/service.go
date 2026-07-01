// Package service implements category read logic for the public API.
package service

// 本文件实现分类读取的业务逻辑，保持 handler -> service -> repository 单向依赖。

import (
	"context"

	categoryrepo "gotribe/internal/api/category/repository"
	"gotribe/internal/model"
)

// Service 负责封装分类相关的业务逻辑。
type Service struct {
	repo *categoryrepo.Repository
}

// NewService 创建分类服务实例。
func NewService(repo *categoryrepo.Repository) *Service {
	return &Service{repo: repo}
}

// ListByParentID 返回指定父分类下的启用子分类。
func (s *Service) ListByParentID(ctx context.Context, parentID int64) ([]model.Category, error) {
	return s.repo.ListByParentID(ctx, parentID)
}
