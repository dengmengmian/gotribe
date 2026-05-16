package service

import (
	"context"

	"gotribe/internal/admin/menu/dto"
	"gotribe/internal/admin/menu/repository"
	"gotribe/internal/core/database"
	"gotribe/internal/model"
)

// Service 菜单业务逻辑接口
type Service interface {
	List(ctx context.Context) ([]*model.Menu, error)
	Tree(ctx context.Context) ([]*model.Menu, error)
	Create(ctx context.Context, actor model.Admin, req *dto.CreateMenuRequest) error
	Update(ctx context.Context, actor model.Admin, menuID int64, req *dto.UpdateMenuRequest) error
	Delete(ctx context.Context, menuIds []int64) error
	GetUserMenusByUserID(ctx context.Context, userID int64) ([]*model.Menu, error)
	GetUserMenuTreeByUserID(ctx context.Context, userID int64) ([]*model.Menu, error)
}

// service 菜单业务逻辑实现
type service struct {
	menuRepo *repository.Repository
}

// NewMenuService 创建菜单服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		menuRepo: repository.NewRepository(tx),
	}
}

// List 获取菜单列表
func (s *service) List(ctx context.Context) ([]*model.Menu, error) {
	return s.menuRepo.List(ctx)
}

// Tree 获取菜单树
func (s *service) Tree(ctx context.Context) ([]*model.Menu, error) {
	return s.menuRepo.Tree(ctx)
}

// Create 创建菜单
func (s *service) Create(ctx context.Context, actor model.Admin, req *dto.CreateMenuRequest) error {
	menu := model.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       &req.Icon,
		Path:       req.Path,
		Redirect:   &req.Redirect,
		Component:  req.Component,
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
		NoCache:    req.NoCache,
		AlwaysShow: req.AlwaysShow,
		Breadcrumb: req.Breadcrumb,
		ActiveMenu: &req.ActiveMenu,
		ParentID:   &req.ParentID,
		Creator:    actor.Username,
	}
	return s.menuRepo.Create(ctx, &menu)
}

// Update 更新菜单
func (s *service) Update(ctx context.Context, actor model.Admin, menuID int64, req *dto.UpdateMenuRequest) error {
	menu := model.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       &req.Icon,
		Path:       req.Path,
		Redirect:   &req.Redirect,
		Component:  req.Component,
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
		NoCache:    req.NoCache,
		AlwaysShow: req.AlwaysShow,
		Breadcrumb: req.Breadcrumb,
		ActiveMenu: &req.ActiveMenu,
		ParentID:   &req.ParentID,
		Creator:    actor.Username,
	}
	return s.menuRepo.Update(ctx, menuID, &menu)
}

// Delete 批量删除菜单
func (s *service) Delete(ctx context.Context, menuIds []int64) error {
	return s.menuRepo.Delete(ctx, menuIds)
}

// GetUserMenusByUserID 获取用户的可访问菜单列表
func (s *service) GetUserMenusByUserID(ctx context.Context, userID int64) ([]*model.Menu, error) {
	return s.menuRepo.GetUserMenusByUserID(ctx, userID)
}

// GetUserMenuTreeByUserID 获取用户的可访问菜单树
func (s *service) GetUserMenuTreeByUserID(ctx context.Context, userID int64) ([]*model.Menu, error) {
	return s.menuRepo.GetUserMenuTreeByUserID(ctx, userID)
}
