package repository

import (
	"context"
	"fmt"
	"strings"

	"gotribe/internal/admin/config/dto"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
	"gotribe/internal/core/util"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建配置仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// GetConfigByID 获取单个配置
func (r *Repository) GetConfigByID(ctx context.Context, id int64) (model.Config, error) {
	var config model.Config
	err := r.tx.DB(ctx).Where("id = ?", id).First(&config).Error
	return config, err
}

// GetConfigByAlias 根据别名获取单个配置
func (r *Repository) GetConfigByAlias(ctx context.Context, alias string) (model.Config, error) {
	var config model.Config
	err := r.tx.DB(ctx).Where("alias = ?", strings.TrimSpace(alias)).First(&config).Error
	return config, err
}

// List 获取配置列表
func (r *Repository) List(ctx context.Context, req *dto.ConfigListRequest) ([]*model.Config, int64, error) {
	var list []*model.Config
	db := r.tx.DB(ctx).Model(&model.Config{})

	title := strings.TrimSpace(req.Title)
	if !utils.IsEmpty(title) {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	}
	alias := strings.TrimSpace(req.Alias)
	if !utils.IsEmpty(alias) {
		db = db.Where("alias = ?", alias)
	}
	if req.ProjectID > 0 {
		db = db.Where("project_id = ?", req.ProjectID)
	}
	reqType := req.Type
	if reqType != 0 {
		db = db.Where("type = ?", reqType)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	db = db.Order("created_at DESC")

	page, perPage := database.NormalizePagination(int(req.PageNum), int(req.PageSize))
	err = db.Offset((page - 1) * perPage).Limit(perPage).Find(&list).Error

	return r.GetConfigOther(ctx, list), total, err
}

// GetConfigOther 获取配置其他信息
func (r *Repository) GetConfigOther(ctx context.Context, configs []*model.Config) []*model.Config {
	if len(configs) == 0 {
		return configs
	}

	// 批量收集并去重 project_id，避免 N+1 查询
	projectIDSet := make(map[int64]struct{}, len(configs))
	projectIDs := make([]int64, 0, len(configs))
	for _, cfg := range configs {
		if cfg.ProjectID == 0 {
			continue
		}
		if _, exists := projectIDSet[cfg.ProjectID]; exists {
			continue
		}
		projectIDSet[cfg.ProjectID] = struct{}{}
		projectIDs = append(projectIDs, cfg.ProjectID)
	}

	if len(projectIDs) == 0 {
		return configs
	}

	var projects []model.Project
	if err := r.tx.DB(ctx).Where("id IN (?)", projectIDs).Find(&projects).Error; err != nil {
		return configs
	}

	projectMap := make(map[int64]*model.Project, len(projects))
	for i := range projects {
		projectMap[projects[i].ID] = &projects[i]
	}

	for _, m := range configs {
		if project, ok := projectMap[m.ProjectID]; ok {
			m.Project = project
		}
	}
	return configs
}

// Create 创建配置
func (r *Repository) Create(ctx context.Context, config *model.Config) error {
	err := r.tx.DB(ctx).Create(config).Error
	return err
}

// UpdateConfig 更新配置
func (r *Repository) UpdateConfig(ctx context.Context, config *model.Config) error {
	err := r.tx.DB(ctx).Model(config).Updates(config).Error
	if err != nil {
		return err
	}

	return err
}

// Delete 批量删除
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	var configs []model.Config
	for _, id := range ids {
		config, err := r.GetConfigByID(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的配置", id)
		}
		configs = append(configs, config)
	}

	err := r.tx.DB(ctx).Unscoped().Delete(&configs).Error

	return err
}
