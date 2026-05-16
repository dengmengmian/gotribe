package service

import (
	"context"
	"errors"
	"fmt"

	adminrepo "gotribe/internal/admin/admin_user/repository"
	apirepo "gotribe/internal/admin/api/repository"
	"gotribe/internal/core/errs"
	menurepo "gotribe/internal/admin/menu/repository"
	"gotribe/internal/admin/role/dto"
	"gotribe/internal/admin/role/repository"
	"gotribe/internal/core/database"
	"gotribe/internal/model"

	"github.com/casbin/casbin/v2"
	"github.com/thoas/go-funk"
)

// Service 角色业务逻辑接口
type Service interface {
	List(ctx context.Context, req *dto.RoleListRequest) ([]dto.RoleResponse, int64, error)
	Create(ctx context.Context, actor model.Admin, req *dto.CreateRoleRequest) error
	Update(ctx context.Context, actor model.Admin, roleID int64, req *dto.CreateRoleRequest) error
	GetRoleMenusByID(ctx context.Context, roleID int64) ([]dto.MenuSummary, error)
	UpdateRoleMenusByID(ctx context.Context, actor model.Admin, roleID int64, req *dto.UpdateRoleMenusRequest) error
	GetRoleApisByID(ctx context.Context, roleID int64) ([]dto.ApiSummary, error)
	UpdateRoleApisByID(ctx context.Context, actor model.Admin, roleID int64, req *dto.UpdateRoleApisRequest) error
	Delete(ctx context.Context, actor model.Admin, roleIds []int64) error
}

// service 角色业务逻辑实现
type service struct {
	roleRepo  *repository.Repository
	adminRepo *adminrepo.Repository
	menuRepo  *menurepo.Repository
	apiRepo   *apirepo.Repository
	enforcer  *casbin.Enforcer
	tx        *database.TransactionManager
}

// NewRoleService 创建角色服务实例
func NewService(tx *database.TransactionManager, enforcer *casbin.Enforcer) Service {
	return &service{
		roleRepo:  repository.NewRepository(tx, enforcer),
		adminRepo: adminrepo.NewRepository(tx),
		menuRepo:  menurepo.NewRepository(tx),
		apiRepo:   apirepo.NewRepository(tx, enforcer),
		enforcer:  enforcer,
		tx:        tx,
	}
}

// List 获取角色列表。
func (s *service) List(ctx context.Context, req *dto.RoleListRequest) ([]dto.RoleResponse, int64, error) {
	roles, total, err := s.roleRepo.List(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	return dto.ToRoleListResponse(roles), total, nil
}

// CreateRole 创建角色
func (s *service) Create(ctx context.Context, actor model.Admin, req *dto.CreateRoleRequest) error {
	sort, ctxUser, err := s.adminRepo.GetCurrentAdminMinRoleSort(ctx, actor)
	if err != nil {
		return err
	}

	if sort >= int64(req.Sort) {
		return errors.New("不能创建比自己等级高或相同等级的角色")
	}

	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}

	return s.roleRepo.Create(ctx, &role)
}

// UpdateRoleByID 更新角色
func (s *service) Update(ctx context.Context, actor model.Admin, roleID int64, req *dto.CreateRoleRequest) error {
	minSort, ctxUser, err := s.adminRepo.GetCurrentAdminMinRoleSort(ctx, actor)
	if err != nil {
		return err
	}

	roles, err := s.roleRepo.GetRolesByIds(ctx, []int64{roleID})
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return errors.New("未获取到角色信息")
	}
	targetRole := roles[0]
	if uint64(targetRole.Sort) < uint64(minSort) {
		return errors.New("不能更新比自己角色等级高的角色")
	}

	if uint64(targetRole.Sort) == uint64(minSort) {
		req.Keyword = targetRole.Keyword
		req.Sort = targetRole.Sort
		req.Status = targetRole.Status
	} else {
		if minSort >= int64(req.Sort) {
			return errors.New("不能把角色等级更新得比当前用户的等级高或相同")
		}
	}

	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}

	err = s.roleRepo.Update(ctx, roleID, &role)
	if err != nil {
		return err
	}

	if req.Keyword != roles[0].Keyword {
		rolePolicies, _ := s.enforcer.GetFilteredPolicy(0, roles[0].Keyword)
		if len(rolePolicies) == 0 {
			s.adminRepo.ClearAdminInfoCache()
			return nil
		}
		rolePoliciesCopy := make([][]string, 0)
		for _, policy := range rolePolicies {
			policyCopy := make([]string, len(policy))
			copy(policyCopy, policy)
			rolePoliciesCopy = append(rolePoliciesCopy, policyCopy)
			policy[0] = req.Keyword
		}

		isAdded, _ := s.enforcer.AddPolicies(rolePolicies)
		if !isAdded {
			return errors.New("更新角色成功，但角色关键字关联的权限接口更新失败")
		}
		isRemoved, _ := s.enforcer.RemovePolicies(rolePoliciesCopy)
		if !isRemoved {
			return errors.New("更新角色成功，但角色关键字关联的权限接口更新失败")
		}
		err := s.enforcer.LoadPolicy()
		if err != nil {
			return errors.New("更新角色成功，但角色关键字关联角色的权限接口策略加载失败")
		}
	}

	s.adminRepo.ClearAdminInfoCache()
	return nil
}

// GetRoleMenusByID 获取角色的权限菜单。
func (s *service) GetRoleMenusByID(ctx context.Context, roleID int64) ([]dto.MenuSummary, error) {
	menus, err := s.roleRepo.GetRoleMenusByID(ctx, roleID)
	if err != nil {
		return nil, err
	}
	return dto.ToMenuSummaryList(menus), nil
}

// UpdateRoleMenusByID 更新角色的权限菜单
func (s *service) UpdateRoleMenusByID(ctx context.Context, actor model.Admin, roleID int64, req *dto.UpdateRoleMenusRequest) error {
	roles, err := s.roleRepo.GetRolesByIds(ctx, []int64{roleID})
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return errors.New("未获取到角色信息")
	}

	minSort, ctxUser, err := s.adminRepo.GetCurrentAdminMinRoleSort(ctx, actor)
	if err != nil {
		return err
	}

	if minSort != 1 {
		if minSort >= int64(roles[0].Sort) {
			return errors.New("不能更新比自己角色等级高或相等角色的权限菜单")
		}
	}

	ctxUserMenus, err := s.menuRepo.GetUserMenusByUserID(ctx, ctxUser.ID)
	if err != nil {
		return err
	}

	ctxUserMenusIds := make([]int64, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.ID)
	}

	menuIds := req.MenuIds
	reqMenus := make([]*model.Menu, 0)

	if minSort != 1 {
		for _, id := range menuIds {
			if !funk.Contains(ctxUserMenusIds, id) {
				return errors.New(fmt.Sprintf("无权设置ID为%d的菜单", id))
			}
		}

		for _, id := range menuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.ID {
					reqMenus = append(reqMenus, menu)
					break
				}
			}
		}
	} else {
		menus, err := s.menuRepo.List(ctx)
		if err != nil {
			return err
		}
		for _, menuID := range menuIds {
			for _, menu := range menus {
				if menuID == menu.ID {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}

	roles[0].Menus = reqMenus
	return s.roleRepo.UpdateRoleMenus(ctx, roles[0])
}

// GetRoleApisByID 获取角色的权限接口。
func (s *service) GetRoleApisByID(ctx context.Context, roleID int64) ([]dto.ApiSummary, error) {
	roles, err := s.roleRepo.GetRolesByIds(ctx, []int64{roleID})
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, errors.New("未获取到角色信息")
	}
	apis, err := s.roleRepo.GetRoleApisByRoleKeyword(ctx, roles[0].Keyword)
	if err != nil {
		return nil, err
	}
	return dto.ToApiSummaryList(apis), nil
}

// UpdateRoleApisByID 更新角色的权限接口
func (s *service) UpdateRoleApisByID(ctx context.Context, actor model.Admin, roleID int64, req *dto.UpdateRoleApisRequest) error {
	roles, err := s.roleRepo.GetRolesByIds(ctx, []int64{roleID})
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return errors.New("未获取到角色信息")
	}

	minSort, ctxUser, err := s.adminRepo.GetCurrentAdminMinRoleSort(ctx, actor)
	if err != nil {
		return err
	}

	if minSort != 1 {
		if minSort >= int64(roles[0].Sort) {
			return errors.New("不能更新比自己角色等级高或相等角色的权限接口")
		}
	}

	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy, _ := s.enforcer.GetFilteredPolicy(0, role.Keyword)
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}

	apiIds := req.ApiIds
	apis, err := s.apiRepo.GetApisByID(ctx, apiIds)
	if err != nil {
		return errors.New("根据接口ID获取接口信息失败")
	}

	reqRolePolicies := make([][]string, 0)
	for _, api := range apis {
		reqRolePolicies = append(reqRolePolicies, []string{
			roles[0].Keyword, api.Path, api.Method,
		})
	}

	if minSort != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				return errors.New(fmt.Sprintf("无权设置路径为%s,请求方式为%s的接口", reqPolicy[1], reqPolicy[2]))
			}
		}
	}

	return s.roleRepo.UpdateRoleApis(ctx, roles[0].Keyword, reqRolePolicies)
}

// BatchDeleteRoleByIds 批量删除角色
func (s *service) Delete(ctx context.Context, actor model.Admin, roleIds []int64) error {
	minSort, _, err := s.adminRepo.GetCurrentAdminMinRoleSort(ctx, actor)
	if err != nil {
		return err
	}

	roles, err := s.roleRepo.GetRolesByIds(ctx, roleIds)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return errors.New("未获取到角色信息")
	}

	for _, role := range roles {
		if minSort >= int64(role.Sort) {
			return errors.New("不能删除比自己角色等级高或相等的角色")
		}
	}

	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.roleRepo.Delete(txCtx, roleIds)
	})
	if err != nil {
		return err
	}

	for _, role := range roles {
		rmPolicies, _ := s.enforcer.GetFilteredPolicy(0, role.Keyword)
		if len(rmPolicies) > 0 {
			isRemoved, _ := s.enforcer.RemovePolicies(rmPolicies)
			if !isRemoved {
				return errs.InternalWithKey(errs.MsgDeleteRoleApisFailed, nil, nil)
			}
		}
	}

	s.adminRepo.ClearAdminInfoCache()
	return nil
}
