package repository

import (
	"context"
	"fmt"
	"strings"

	"gotribe/internal/admin/role/dto"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	"gotribe/internal/model"

	"github.com/casbin/casbin/v2"
)

// Repository 角色数据访问实现
type Repository struct {
	tx       *database.TransactionManager
	enforcer *casbin.SyncedEnforcer
}

// NewRepository 创建角色仓库实例
func NewRepository(tx *database.TransactionManager, enforcer *casbin.SyncedEnforcer) *Repository {
	return &Repository{tx: tx, enforcer: enforcer}
}

func buildRoleOrder(req *dto.RoleListRequest) string {
	sortByMap := map[string]string{
		"name":       "name",
		"keyword":    "keyword",
		"sort":       "sort",
		"status":     "status",
		"creator":    "creator",
		"createdAt":  "created_at",
		"created_at": "created_at",
	}

	column, ok := sortByMap[strings.TrimSpace(req.SortBy)]
	if !ok {
		return "created_at DESC"
	}

	direction := "ASC"
	if strings.EqualFold(strings.TrimSpace(req.SortOrder), "desc") {
		direction = "DESC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

// List 获取角色列表
func (r *Repository) List(ctx context.Context, req *dto.RoleListRequest) ([]model.Role, int64, error) {
	var list []model.Role
	db := r.tx.DB(ctx).Model(&model.Role{})

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	db = db.Order(buildRoleOrder(req))

	page, perPage := database.NormalizePagination(int(req.PageNum), int(req.PageSize))
	err = db.Offset((page - 1) * perPage).Limit(perPage).Find(&list).Error

	return list, total, err
}

// GetRolesByIds 根据角色ID获取角色
func (r *Repository) GetRolesByIds(ctx context.Context, roleIds []int64) ([]*model.Role, error) {
	var list []*model.Role
	err := r.tx.DB(ctx).Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// Create 创建角色
func (r *Repository) Create(ctx context.Context, role *model.Role) error {
	err := r.tx.DB(ctx).Create(role).Error
	return err
}

// Update 更新角色
func (r *Repository) Update(ctx context.Context, roleID int64, role *model.Role) error {
	err := r.tx.DB(ctx).Model(&model.Role{}).Where("id = ?", roleID).Updates(role).Error
	return err
}

// GetRoleMenusByID 获取角色的权限菜单
func (r *Repository) GetRoleMenusByID(ctx context.Context, roleID int64) ([]*model.Menu, error) {
	var role model.Role
	err := r.tx.DB(ctx).Where("id = ?", roleID).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// UpdateRoleMenus 更新角色的权限菜单
func (r *Repository) UpdateRoleMenus(ctx context.Context, role *model.Role) error {
	err := r.tx.DB(ctx).Model(role).Association("Menus").Replace(role.Menus)
	return err
}

// GetRoleApisByRoleKeyword 根据角色关键字获取角色的权限接口
func (r *Repository) GetRoleApisByRoleKeyword(ctx context.Context, roleKeyword string) ([]*model.Api, error) {
	policies, _ := r.enforcer.GetFilteredPolicy(0, roleKeyword)

	// 获取所有接口
	var apis []*model.Api
	err := r.tx.DB(ctx).Find(&apis).Error
	if err != nil {
		return apis, errs.InternalWithKey(errs.MsgGetRoleApisFailed, nil, nil)
	}

	accessApis := make([]*model.Api, 0)

	for _, policy := range policies {
		path := policy[1]
		method := policy[2]
		for _, api := range apis {
			if path == api.Path && method == api.Method {
				accessApis = append(accessApis, api)
				break
			}
		}
	}

	return accessApis, err

}

// UpdateRoleApis 更新角色的权限接口（先全部删除再新增）
func (r *Repository) UpdateRoleApis(ctx context.Context, roleKeyword string, reqRolePolicies [][]string) error {
	// 先获取path中的角色ID对应角色已有的police(需要先删除的)
	err := r.enforcer.LoadPolicy()
	if err != nil {
		return errs.InternalWithKey(errs.MsgLoadRolePolicyFailed, nil, nil)
	}
	rmPolicies, _ := r.enforcer.GetFilteredPolicy(0, roleKeyword)
	if len(rmPolicies) > 0 {
		isRemoved, err := r.enforcer.RemovePolicies(rmPolicies)
		if err != nil {
			return errs.InternalWithKey(errs.MsgUpdateRoleApisFailed, nil, err)
		}
		if !isRemoved {
			return errs.InternalWithKey(errs.MsgUpdateRoleApisFailed, nil, nil)
		}
	}
	isAdded, err := r.enforcer.AddPolicies(reqRolePolicies)
	if err != nil {
		return errs.InternalWithKey(errs.MsgUpdateRoleApisFailed, nil, err)
	}
	if !isAdded {
		return errs.InternalWithKey(errs.MsgUpdateRoleApisFailed, nil, nil)
	}
	err = r.enforcer.LoadPolicy()
	if err != nil {
		return errs.InternalWithKey(errs.MsgLoadRolePolicyFailed, nil, nil)
	} else {
		return err
	}
}

// Delete 删除角色
func (r *Repository) Delete(ctx context.Context, roleIds []int64) error {
	var roles []*model.Role
	db := r.tx.DB(ctx)
	if err := db.Where("id IN (?)", roleIds).Find(&roles).Error; err != nil {
		return err
	}
	if err := db.Select("Users", "Menus").Unscoped().Delete(&roles).Error; err != nil {
		return err
	}
	return nil
}
