// Package service implements tag read logic for the public API.
package service

// 本文件实现标签读取的业务逻辑，保持 handler -> service -> repository 单向依赖。

import (
	"context"

	tagrepo "gotribe/internal/api/tag/repository"
	tagmodel "gotribe/internal/model"
)

// maxTagPerPage 标签列表单页上限，防止客户端传超大 per_page 造成无界查询。
const maxTagPerPage = 100

// Service 负责封装标签相关的业务逻辑。
type Service struct {
	repo *tagrepo.Repository
}

// NewService 创建标签服务实例。
func NewService(repo *tagrepo.Repository) *Service {
	return &Service{repo: repo}
}

// List 返回启用标签，perPage 归一到 (0, maxTagPerPage]。
func (s *Service) List(ctx context.Context, keyword string, perPage int) ([]tagmodel.Tag, error) {
	if perPage <= 0 || perPage > maxTagPerPage {
		perPage = maxTagPerPage
	}
	return s.repo.List(ctx, keyword, perPage)
}
