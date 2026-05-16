package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gotribe/internal/admin/api/dto"
	"gotribe/internal/core/database"
	"gotribe/internal/model"

	"github.com/casbin/casbin/v2"
	"github.com/thoas/go-funk"
)

type Repository struct {
	tx       *database.TransactionManager
	enforcer *casbin.Enforcer
}

func NewRepository(tx *database.TransactionManager, enforcer *casbin.Enforcer) *Repository {
	return &Repository{tx: tx, enforcer: enforcer}
}

func buildApiOrder(req *dto.ApiListRequest) string {
	sortByMap := map[string]string{
		"method":     "method",
		"path":       "path",
		"category":   "category",
		"desc":       "\"desc\"",
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

// List 获取接口列表
func (a *Repository) List(ctx context.Context, req *dto.ApiListRequest) ([]*model.Api, int64, error) {
	var list []*model.Api
	db := a.tx.DB(ctx).Model(&model.Api{})

	method := strings.TrimSpace(req.Method)
	if method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		db = db.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	// 当page > 0 且 perPage > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	db = db.Order(buildApiOrder(req))
	page := int(req.PageNum)
	perPage := int(req.PageSize)
	if page > 0 && perPage > 0 {
		err = db.Offset((page - 1) * perPage).Limit(perPage).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// GetApisByID 根据接口ID获取接口列表
func (a *Repository) GetApisByID(ctx context.Context, apiIds []int64) ([]*model.Api, error) {
	var apis []*model.Api
	err := a.tx.DB(ctx).Where("id IN (?)", apiIds).Find(&apis).Error
	return apis, err
}

// Tree 获取接口树(按接口Category字段分类)
func (a *Repository) Tree(ctx context.Context) ([]*dto.ApiTreeResponse, error) {
	var apiList []*model.Api
	err := a.tx.DB(ctx).Order("category").Order("created_at").Find(&apiList).Error
	// 获取所有的分类
	var categoryList []string
	for _, api := range apiList {
		categoryList = append(categoryList, api.Category)
	}
	// 获取去重后的分类
	categoryUniq := funk.UniqString(categoryList)

	apiTree := make([]*dto.ApiTreeResponse, len(categoryUniq))

	for i, category := range categoryUniq {
		apiTree[i] = &dto.ApiTreeResponse{
			ID:       -i,
			Desc:     category,
			Category: category,
			Children: nil,
		}
		for _, api := range apiList {
			if category == api.Category {
				apiTree[i].Children = append(apiTree[i].Children, api)
			}
		}
	}

	return apiTree, err
}

// Create 创建接口
func (a *Repository) Create(ctx context.Context, api *model.Api) error {
	err := a.tx.DB(ctx).Create(api).Error
	return err
}

// Update 更新接口
func (a *Repository) Update(ctx context.Context, apiID int64, api *model.Api) error {
	// 根据id获取接口信息
	var oldApi model.Api
	err := a.tx.DB(ctx).First(&oldApi, apiID).Error
	if err != nil {
		return errors.New("根据接口ID获取接口信息失败")
	}
	err = a.tx.DB(ctx).Model(api).Where("id = ?", apiID).Updates(api).Error
	if err != nil {
		return err
	}
	// 更新了method和path就更新casbin中policy
	if oldApi.Path != api.Path || oldApi.Method != api.Method {
		policies, _ := a.enforcer.GetFilteredPolicy(1, oldApi.Path, oldApi.Method)
		// 接口在casbin的policy中存在才进行操作
		if len(policies) > 0 {
			// 先删除
			isRemoved, _ := a.enforcer.RemovePolicies(policies)
			if !isRemoved {
				return errors.New("更新权限接口失败")
			}
			for _, policy := range policies {
				policy[1] = api.Path
				policy[2] = api.Method
			}
			// 新增
			isAdded, _ := a.enforcer.AddPolicies(policies)
			if !isAdded {
				return errors.New("更新权限接口失败")
			}
			// 加载policy
			err := a.enforcer.LoadPolicy()
			if err != nil {
				return errors.New("更新权限接口成功，权限接口策略加载失败")
			} else {
				return err
			}
		}
	}
	return err
}

// Delete 批量删除接口
func (a *Repository) Delete(ctx context.Context, apiIds []int64) error {
	apis, err := a.GetApisByID(ctx, apiIds)
	if err != nil {
		return errors.New("根据接口ID获取接口列表失败")
	}
	if len(apis) == 0 {
		return errors.New("根据接口ID未获取到接口列表")
	}

	err = a.tx.DB(ctx).Where("id IN (?)", apiIds).Unscoped().Delete(&model.Api{}).Error
	// 如果删除成功，删除casbin中policy
	if err == nil {
		for _, api := range apis {
			policies, _ := a.enforcer.GetFilteredPolicy(1, api.Path, api.Method)
			if len(policies) > 0 {
				isRemoved, _ := a.enforcer.RemovePolicies(policies)
				if !isRemoved {
					return errors.New("删除权限接口失败")
				}
			}
		}
		// 重新加载策略
		err := a.enforcer.LoadPolicy()
		if err != nil {
			return errors.New("删除权限接口成功，权限接口策略加载失败")
		} else {
			return err
		}
	}
	return err
}

// GetApiDescByPath 根据接口路径和请求方式获取接口描述
func (a *Repository) GetApiDescByPath(ctx context.Context, path string, method string) (string, error) {
	var api model.Api
	err := a.tx.DB(ctx).Where("path = ?", path).Where("method = ?", method).First(&api).Error
	return api.Desc, err
}
