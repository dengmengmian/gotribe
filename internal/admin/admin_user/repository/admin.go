package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gotribe/internal/admin/admin_user/dto"
	"gotribe/internal/core/constant"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/util"
	"gotribe/internal/model"

	"github.com/patrickmn/go-cache"
	"github.com/thoas/go-funk"
)

// Repository 管理员数据访问实现
type Repository struct {
	tx *database.TransactionManager
}

// 当前用户信息缓存，避免频繁获取数据库。
// TTL 取较短值：管理员被禁用 / 角色变更等安全敏感变更最长在此时间内生效
// （多数变更路径会显式失效缓存，此 TTL 是兜底上界）。
var adminInfoCache = cache.New(5*time.Minute, 10*time.Minute)

// dummyAdminHash 用于时序攻击防护：无论管理员是否存在都执行一次 bcrypt 比较，
// 使响应耗时一致，避免用户名枚举。值为一个合法但无意义的 bcrypt 哈希。
const dummyAdminHash = "$2a$10$abcdefghijklmnopqrstuuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

// NewRepository 创建管理员仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

func buildAdminOrder(req *dto.AdminListRequest) string {
	sortByMap := map[string]string{
		"username":   "username",
		"nickname":   "nickname",
		"status":     "status",
		"mobile":     "mobile",
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

// Login 登录
func (r *Repository) Login(ctx context.Context, admin *model.Admin) (*model.Admin, error) {
	// 根据用户名获取用户
	var firstAdmin model.Admin
	err := r.tx.DB(ctx).
		Where("username = ?", admin.Username).
		Preload("Roles").
		First(&firstAdmin).Error
	if err != nil {
		// 用户不存在也执行一次 bcrypt 比较，消除时序差异，避免用户名枚举；
		// 且返回与「密码错误」相同的通用消息，不泄露账户是否存在。
		_ = utils.PasswordUtil.ComparePasswd(dummyAdminHash, admin.Password)
		return nil, errs.Unauthorized(errs.T("zh", errs.MsgInvalidCredentials))
	}

	// 先校验密码：在攻击者证明知道密码之前，不泄露账户是否被禁用/角色状态（防枚举）。
	if err := utils.PasswordUtil.ComparePasswd(firstAdmin.Password, admin.Password); err != nil {
		return nil, errs.Unauthorized(errs.T("zh", errs.MsgInvalidCredentials))
	}

	// 密码正确后再检查账户与角色状态。
	if firstAdmin.Status != 1 {
		return nil, errs.Forbidden(errs.T("zh", errs.MsgUserDisabled))
	}

	roles := firstAdmin.Roles
	isValidate := false
	for _, role := range roles {
		// 有一个正常状态的角色就可以登录
		if role.Status == constant.DEFAULT_ID {
			isValidate = true
			break
		}
	}
	if !isValidate {
		return nil, errs.Forbidden(errs.T("zh", errs.MsgUserRoleDisabled))
	}

	return &firstAdmin, nil
}

// Me 获取当前登录用户信息
// 需要缓存，减少数据库访问
func (r *Repository) Me(ctx context.Context, actor model.Admin) (model.Admin, error) {
	var newAdmin model.Admin
	if actor.ID == 0 || actor.Username == "" {
		return newAdmin, errs.Unauthorized(errs.T("zh", errs.MsgUserNotLoggedIn))
	}

	// 先获取缓存
	cacheAdmin, found := adminInfoCache.Get(actor.Username)
	var admin model.Admin
	var err error
	if found {
		admin = cacheAdmin.(model.Admin)
		err = nil
	} else {
		// 缓存中没有就获取数据库
		admin, err = r.GetAdminByID(ctx, actor.ID)
		// 获取成功就缓存
		if err != nil {
			adminInfoCache.Delete(actor.Username)
		} else {
			adminInfoCache.Set(actor.Username, admin, cache.DefaultExpiration)
		}
	}
	return admin, err
}

// GetCurrentAdminMinRoleSort 获取当前用户角色排序最小值（最高等级角色）以及当前用户信息
func (r *Repository) GetCurrentAdminMinRoleSort(ctx context.Context, actor model.Admin) (int64, model.Admin, error) {
	// 获取当前用户
	ctxAdmin, err := r.Me(ctx, actor)
	if err != nil {
		return 999, ctxAdmin, err
	}
	// 获取当前用户的所有角色
	currentRoles := ctxAdmin.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := int64(funk.MinInt(currentRoleSorts))

	return currentRoleSortMin, ctxAdmin, nil
}

// GetAdminByID 获取单个用户
func (r *Repository) GetAdminByID(ctx context.Context, id int64) (model.Admin, error) {
	var admin model.Admin
	err := r.tx.DB(ctx).Where("id = ?", id).Preload("Roles").First(&admin).Error
	return admin, err
}

// List 获取用户列表
func (r *Repository) List(ctx context.Context, req *dto.AdminListRequest) ([]*model.Admin, int64, error) {
	var list []*model.Admin
	db := r.tx.DB(ctx).Model(&model.Admin{})

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当page > 0 且 perPage > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	db = db.Order(buildAdminOrder(req))
	page := int(req.PageNum)
	perPage := int(req.PageSize)
	if page > 0 && perPage > 0 {
		err = db.Offset((page - 1) * perPage).Limit(perPage).Preload("Roles").Find(&list).Error
	} else {
		err = db.Preload("Roles").Find(&list).Error
	}
	return list, total, err
}

// UpdatePassword 更新密码
func (r *Repository) UpdatePassword(ctx context.Context, username string, hashNewPasswd string) error {
	err := r.tx.DB(ctx).Model(&model.Admin{}).Where("username = ?", username).Update("password", hashNewPasswd).Error
	// 如果更新密码成功，则更新当前用户信息缓存
	// 先获取缓存
	cacheAdmin, found := adminInfoCache.Get(username)
	if err == nil {
		if found {
			admin := cacheAdmin.(model.Admin)
			admin.Password = hashNewPasswd
			adminInfoCache.Set(username, admin, cache.DefaultExpiration)
		} else {
			// 没有缓存就获取用户信息缓存
			var admin model.Admin
			if err := r.tx.DB(ctx).Where("username = ?", username).First(&admin).Error; err != nil {
				return err
			}
			adminInfoCache.Set(username, admin, cache.DefaultExpiration)
		}
	}

	return err
}

// Create 创建用户
func (r *Repository) Create(ctx context.Context, admin *model.Admin) error {
	err := r.tx.DB(ctx).Create(admin).Error
	return err
}

// UpdateAdmin 更新用户
func (r *Repository) UpdateAdmin(ctx context.Context, admin *model.Admin) error {
	db := r.tx.DB(ctx)
	err := db.Model(admin).Updates(admin).Error
	if err == nil {
		err = db.Model(admin).Association("Roles").Replace(admin.Roles)
	}

	// 如果更新成功就更新用户信息缓存（缓存不在事务内）
	if err == nil {
		adminInfoCache.Set(admin.Username, *admin, cache.DefaultExpiration)
	}
	return err
}

// Delete 批量删除
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	// 用户和角色存在多对多关联关系
	var admins []model.Admin
	if err := r.tx.DB(ctx).Where("id IN ?", ids).Preload("Roles").Find(&admins).Error; err != nil {
		return err
	}
	if len(admins) != len(ids) {
		return errs.NotFoundWithKey(errs.MsgUserNotFoundByID, []any{0}, nil)
	}

	err := r.tx.DB(ctx).Select("Roles").Unscoped().Delete(&admins).Error
	// 删除用户成功，则删除用户信息缓存
	if err == nil {
		for _, admin := range admins {
			adminInfoCache.Delete(admin.Username)
		}
	}
	return err
}

// GetAdminMinRoleSortsByIds 根据用户ID获取用户角色排序最小值
func (r *Repository) GetAdminMinRoleSortsByIds(ctx context.Context, ids []int64) ([]int, error) {
	// 根据用户ID获取用户信息
	var adminList []model.Admin
	err := r.tx.DB(ctx).Where("id IN (?)", ids).Preload("Roles").Find(&adminList).Error
	if err != nil {
		return []int{}, err
	}
	if len(adminList) == 0 {
		return []int{}, errs.NotFoundWithKey(errs.MsgNoUsersFound, nil, nil)
	}
	var roleMinSortList []int
	for _, admin := range adminList {
		roles := admin.Roles
		var roleSortList []int
		for _, role := range roles {
			roleSortList = append(roleSortList, int(role.Sort))
		}
		roleMinSort := funk.MinInt(roleSortList)
		roleMinSortList = append(roleMinSortList, roleMinSort)
	}
	return roleMinSortList, nil
}

// SetAdminInfoCache 设置用户信息缓存
func (r *Repository) SetAdminInfoCache(username string, admin model.Admin) {
	adminInfoCache.Set(username, admin, cache.DefaultExpiration)
}

// UpdateAdminInfoCacheByRoleID 根据角色ID更新拥有该角色的用户信息缓存
func (r *Repository) UpdateAdminInfoCacheByRoleID(ctx context.Context, roleID int64) error {
	var role model.Role
	err := r.tx.DB(ctx).Where("id = ?", roleID).Preload("Admins").First(&role).Error
	if err != nil {
		return errs.InternalWithKey(errs.MsgRoleInfoFailed, nil, nil)
	}

	admins := role.Admin
	if len(admins) == 0 {
		return errs.NotFoundWithKey(errs.MsgNoUsersWithRole, nil, nil)
	}

	// 更新用户信息缓存
	for _, admin := range admins {
		_, found := adminInfoCache.Get(admin.Username)
		if found {
			adminInfoCache.Set(admin.Username, *admin, cache.DefaultExpiration)
		}
	}

	return err
}

// ClearAdminInfoCache 清理所有用户信息缓存
func (r *Repository) ClearAdminInfoCache() {
	adminInfoCache.Flush()
}
