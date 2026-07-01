package service

import (
	"context"
	"errors"

	"gotribe/internal/admin/admin_user/dto"
	"gotribe/internal/admin/admin_user/repository"
	"gotribe/internal/core/errs"
	rolerepo "gotribe/internal/admin/role/repository"
	"gotribe/internal/core/database"
	"gotribe/internal/model"
	"gotribe/internal/core/util"

	"github.com/casbin/casbin/v2"
	"github.com/thoas/go-funk"
)

// Service 管理员业务逻辑接口
type Service interface {
	Me(ctx context.Context, actor model.Admin) (model.Admin, error)
	List(ctx context.Context, req *dto.AdminListRequest) ([]*model.Admin, int64, error)
	UpdatePassword(ctx context.Context, actor model.Admin, oldPassword, newPassword string) error
	Create(ctx context.Context, actor model.Admin, req *dto.CreateAdminRequest) error
	Update(ctx context.Context, actor model.Admin, id int64, req *dto.CreateAdminRequest) error
	Delete(ctx context.Context, actor model.Admin, ids []int64) error
	Detail(ctx context.Context, id int64) (*model.Admin, error)
}

// service 管理员业务逻辑实现
type service struct {
	adminRepo *repository.Repository
	roleRepo  *rolerepo.Repository
	tx        *database.TransactionManager
}

// NewAdminService 创建管理员服务实例
func NewService(tx *database.TransactionManager, enforcer *casbin.SyncedEnforcer) Service {
	return &service{
		adminRepo: repository.NewRepository(tx),
		roleRepo:  rolerepo.NewRepository(tx, enforcer),
		tx:        tx,
	}
}

// GetCurrentAdmin 获取当前登录管理员
func (s *service) Me(ctx context.Context, actor model.Admin) (model.Admin, error) {
	return s.adminRepo.Me(ctx, actor)
}

// GetAdmins 获取管理员列表
func (s *service) List(ctx context.Context, req *dto.AdminListRequest) ([]*model.Admin, int64, error) {
	return s.adminRepo.List(ctx, req)
}

// ChangePwd 修改密码
func (s *service) UpdatePassword(ctx context.Context, actor model.Admin, oldPassword, newPassword string) error {
	user, err := s.adminRepo.Me(ctx, actor)
	if err != nil {
		return err
	}
	correctPasswd := user.Password
	err = utils.PasswordUtil.ComparePasswd(correctPasswd, oldPassword)
	if err != nil {
		return errs.Unauthorized(errs.T("zh", errs.MsgPasswordIncorrect))
	}
	hashedPassword, err := utils.PasswordUtil.GenPasswd(newPassword)
	if err != nil {
		return err
	}
	return s.adminRepo.UpdatePassword(ctx, user.Username, hashedPassword)
}

// CreateAdmin 创建管理员
func (s *service) Create(ctx context.Context, actor model.Admin, req *dto.CreateAdminRequest) error {
	currentRoleSortMin, ctxAdmin, err := s.adminRepo.GetCurrentAdminMinRoleSort(ctx, actor)
	if err != nil {
		return err
	}

	reqRoleIds := req.RoleIds
	roles, err := s.roleRepo.GetRolesByIds(ctx, reqRoleIds)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return errors.New("未获取到角色信息")
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	reqRoleSortMin := int64(funk.MinInt(reqRoleSorts))

	if currentRoleSortMin >= reqRoleSortMin {
		return errors.New("用户不能创建比自己等级高的或者相同等级的用户")
	}

	password := req.Password
	if password == "" {
		// 默认密码，首次登录后请务必修改
		password = "Gotribe!23456"
	}
	hashedPassword, err := utils.PasswordUtil.GenPasswd(password)
	if err != nil {
		return err
	}
	user := model.Admin{
		Username:     req.Username,
		Password:     hashedPassword,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxAdmin.Username,
		Roles:        roles,
	}

	return s.adminRepo.Create(ctx, &user)
}

// UpdateAdminByID 更新管理员
func (s *service) Update(ctx context.Context, actor model.Admin, id int64, req *dto.CreateAdminRequest) error {
	oldAdmin, err := s.adminRepo.GetAdminByID(ctx, id)
	if err != nil {
		return err
	}

	ctxAdmin, err := s.adminRepo.Me(ctx, actor)
	if err != nil {
		return err
	}
	currentRoles := ctxAdmin.Roles
	var currentRoleSorts []int
	var currentRoleIds []int64
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
		currentRoleIds = append(currentRoleIds, role.ID)
	}
	currentRoleSortMin := funk.MinInt(currentRoleSorts)

	reqRoleIds := req.RoleIds
	roles, err := s.roleRepo.GetRolesByIds(ctx, reqRoleIds)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return errors.New("未获取到角色信息")
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	reqRoleSortMin := funk.MinInt(reqRoleSorts)

	user := model.Admin{
		Model:        oldAdmin.Model,
		Username:     req.Username,
		Password:     oldAdmin.Password,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxAdmin.Username,
		Roles:        roles,
	}

	if id == ctxAdmin.ID {
		if req.Status == 2 {
			return errors.New("不能禁用自己")
		}
		reqDiff, currentDiff := funk.Difference(req.RoleIds, currentRoleIds)
		if len(reqDiff.([]int64)) > 0 || len(currentDiff.([]int64)) > 0 {
			return errors.New("不能更改自己的角色")
		}
		if req.Password != "" {
			return errors.New("请到个人中心更新自身密码")
		}
		user.Password = ctxAdmin.Password
	} else {
		minRoleSorts, err := s.adminRepo.GetAdminMinRoleSortsByIds(ctx, []int64{id})
		if err != nil || len(minRoleSorts) == 0 {
			return errors.New("根据用户ID获取用户角色排序最小值失败")
		}
		if currentRoleSortMin >= minRoleSorts[0] {
			return errors.New("用户不能更新比自己角色等级高的或者相同等级的用户")
		}
		if currentRoleSortMin >= reqRoleSortMin {
			return errors.New("用户不能把别的用户角色等级更新得比自己高或相等")
		}
		if req.Password != "" {
			hashedPassword, err := utils.PasswordUtil.GenPasswd(req.Password)
			if err != nil {
				return err
			}
			user.Password = hashedPassword
		}
	}

	return s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.adminRepo.UpdateAdmin(txCtx, &user)
	})
}

// BatchDeleteAdminByIds 批量删除管理员
func (s *service) Delete(ctx context.Context, actor model.Admin, ids []int64) error {
	roleMinSortList, err := s.adminRepo.GetAdminMinRoleSortsByIds(ctx, ids)
	if err != nil || len(roleMinSortList) == 0 {
		return errors.New("根据用户ID获取用户角色排序最小值失败")
	}

	minSort, ctxAdmin, err := s.adminRepo.GetCurrentAdminMinRoleSort(ctx, actor)
	if err != nil {
		return err
	}
	currentRoleSortMin := int(minSort)

	if funk.Contains(ids, ctxAdmin.ID) {
		return errors.New("用户不能删除自己")
	}

	for _, sort := range roleMinSortList {
		if currentRoleSortMin >= sort {
			return errors.New("用户不能删除比自己角色等级高的用户")
		}
	}

	return s.adminRepo.Delete(ctx, ids)
}

// Detail 获取管理员详情。
func (s *service) Detail(ctx context.Context, id int64) (*model.Admin, error) {
	admin, err := s.adminRepo.GetAdminByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
