package repository

import (
	"context"

	"gotribe/internal/core/constant"
	"gotribe/internal/core/database"
	"gotribe/internal/model"

	"github.com/thoas/go-funk"
)

type Repository struct {
	tx *database.TransactionManager
}

func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// List 获取菜单列表
func (m *Repository) List(ctx context.Context) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := m.tx.DB(ctx).Order("sort").Find(&menus).Error
	return menus, err
}

// Tree 获取菜单树
func (m *Repository) Tree(ctx context.Context) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := m.tx.DB(ctx).Order("sort").Find(&menus).Error
	// parentID为0的是根菜单
	return GenMenuTree(0, menus), err
}

func GenMenuTree(parentID int64, menus []*model.Menu) []*model.Menu {
	tree := make([]*model.Menu, 0)
	for _, m := range menus {
		if *m.ParentID == parentID {
			children := GenMenuTree(m.ID, menus)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// Create 创建菜单
func (m *Repository) Create(ctx context.Context, menu *model.Menu) error {
	err := m.tx.DB(ctx).Create(menu).Error
	return err
}

// Update 更新菜单
func (m *Repository) Update(ctx context.Context, menuID int64, menu *model.Menu) error {
	err := m.tx.DB(ctx).Model(menu).Where("id = ?", menuID).Updates(menu).Error
	return err
}

// Delete 批量删除菜单
func (m *Repository) Delete(ctx context.Context, menuIds []int64) error {
	var menus []*model.Menu
	err := m.tx.DB(ctx).Where("id IN (?)", menuIds).Find(&menus).Error
	if err != nil {
		return err
	}
	err = m.tx.DB(ctx).Select("Roles").Unscoped().Delete(&menus).Error
	return err
}

// GetUserMenusByUserID 根据用户ID获取用户的权限(可访问)菜单列表
func (m *Repository) GetUserMenusByUserID(ctx context.Context, userID int64) ([]*model.Menu, error) {
	// 获取用户
	var user model.Admin
	err := m.tx.DB(ctx).Where("id = ?", userID).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	// 获取角色
	roles := user.Roles
	// 所有角色的菜单集合
	roleIDs := make([]int64, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
	}
	var userRoles []model.Role
	if err := m.tx.DB(ctx).Where("id IN ?", roleIDs).Preload("Menus").Find(&userRoles).Error; err != nil {
		return nil, err
	}
	allRoleMenus := make([]*model.Menu, 0)
	for _, userRole := range userRoles {
		allRoleMenus = append(allRoleMenus, userRole.Menus...)
	}

	// 所有角色的菜单集合去重
	allRoleMenusID := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusID = append(allRoleMenusID, int(menu.ID))
	}
	allRoleMenusIDUniq := funk.UniqInt(allRoleMenusID)
	allRoleMenusUniq := make([]*model.Menu, 0)
	for _, id := range allRoleMenusIDUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.ID) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}

	// 获取状态status为1的菜单
	accessMenus := make([]*model.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == constant.DEFAULT_ID {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, err
}

// GetUserMenuTreeByUserID 根据用户ID获取用户的权限(可访问)菜单树
func (m *Repository) GetUserMenuTreeByUserID(ctx context.Context, userID int64) ([]*model.Menu, error) {
	menus, err := m.GetUserMenusByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	tree := GenMenuTree(0, menus)
	return tree, err
}
