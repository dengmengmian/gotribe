package service

import (
	"context"

	"go.uber.org/zap"
	"gotribe/internal/core/database"

	"gotribe/internal/admin/category/dto"
	"gotribe/internal/admin/category/repository"
	"gotribe/internal/model"
)

// ErrSelfParent 分类不能把自己设为父分类
type ErrSelfParent struct{}

func (e ErrSelfParent) Error() string {
	return "不能把自己设为父分类"
}

// Service 分类业务逻辑接口
type Service interface {
	Detail(ctx context.Context, id int64) (model.Category, error)
	List(ctx context.Context) ([]*model.Category, error)
	Tree(ctx context.Context) ([]*model.Category, error)
	Create(ctx context.Context, req *dto.CreateCategoryRequest) error
	Update(ctx context.Context, id int64, req *dto.UpdateCategoryRequest) error
	Delete(ctx context.Context, ids []int64) error
}

// service 分类业务逻辑实现
type service struct {
	categoryRepo *repository.Repository
}

// NewService 创建分类服务实例
func NewService(tx *database.TransactionManager, log *zap.SugaredLogger) Service {
	return &service{
		categoryRepo: repository.NewRepository(tx, log),
	}
}

// Detail 根据ID获取分类
func (s *service) Detail(ctx context.Context, id int64) (model.Category, error) {
	return s.categoryRepo.Detail(ctx, id)
}

// List 获取分类列表
func (s *service) List(ctx context.Context) ([]*model.Category, error) {
	return s.categoryRepo.List(ctx)
}

// Tree 获取分类树
func (s *service) Tree(ctx context.Context) ([]*model.Category, error) {
	return s.categoryRepo.Tree(ctx)
}

// Create 创建分类
func (s *service) Create(ctx context.Context, req *dto.CreateCategoryRequest) error {
	category := model.Category{
		Title:       req.Title,
		Slug:        req.Slug,
		Icon:        req.Icon,
		Path:        req.Path,
		Sort:        req.Sort,
		Status:      req.Status,
		Hidden:      req.Hidden,
		ParentID:    req.ParentID,
		Description: req.Description,
	}
	return s.categoryRepo.Create(ctx, &category)
}

// Update 更新分类
func (s *service) Update(ctx context.Context, id int64, req *dto.UpdateCategoryRequest) error {
	category, err := s.categoryRepo.Detail(ctx, id)
	if err != nil {
		return err
	}
	if req.ParentID == category.ID {
		return ErrSelfParent{}
	}
	category.Title = req.Title
	category.Slug = req.Slug
	category.Icon = req.Icon
	category.Path = req.Path
	category.Sort = req.Sort
	category.Status = req.Status
	category.Hidden = req.Hidden
	category.ParentID = req.ParentID
	category.Description = req.Description
	return s.categoryRepo.Update(ctx, id, &category)
}

// Delete 批量删除分类
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.categoryRepo.Delete(ctx, ids)
}
