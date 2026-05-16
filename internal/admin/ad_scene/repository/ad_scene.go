package repository

import (
	"context"
	"fmt"

	"gotribe/internal/admin/ad_scene/dto"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewAdSceneRepository 创建推广场景仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// Detail 获取单个推广场景
func (r *Repository) Detail(ctx context.Context, id int64) (model.AdScene, error) {
	var adScene model.AdScene
	err := r.tx.DB(ctx).Where("id = ?", id).First(&adScene).Error
	return adScene, err
}

// List 获取推广场景列表
func (r *Repository) List(ctx context.Context, req *dto.AdSceneListRequest) ([]*model.AdScene, int64, error) {
	var list []*model.AdScene
	db := r.tx.DB(ctx).Model(&model.AdScene{})

	if req.ProjectID > 0 {
		db = db.Where("project_id = ?", req.ProjectID)
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

	list, err = r.GetAdSceneOther(ctx, list)
	return list, total, err
}

// GetAdSceneOther 获取推广场景其他信息
func (r *Repository) GetAdSceneOther(ctx context.Context, adScenes []*model.AdScene) ([]*model.AdScene, error) {
	if len(adScenes) == 0 {
		return adScenes, nil
	}

	// 收集所有 ProjectID
	projectIDs := make([]int64, 0, len(adScenes))
	for _, m := range adScenes {
		if m.ProjectID > 0 {
			projectIDs = append(projectIDs, m.ProjectID)
		}
	}

	if len(projectIDs) == 0 {
		return adScenes, nil
	}

	// 批量查询
	var projects []*model.Project
	if err := r.tx.DB(ctx).Where("id IN ?", projectIDs).Find(&projects).Error; err != nil {
		return adScenes, err
	}

	// 建立映射
	projectMap := make(map[int64]*model.Project)
	for _, project := range projects {
		projectMap[project.ID] = project
	}

	// 赋值
	for _, m := range adScenes {
		if project, ok := projectMap[m.ProjectID]; ok {
			m.Project = project
		}
	}
	return adScenes, nil
}

// Create 创建推广场景
func (r *Repository) Create(ctx context.Context, adScene *model.AdScene) error {
	err := r.tx.DB(ctx).Create(adScene).Error
	return err
}

// Update 更新推广场景
func (r *Repository) Update(ctx context.Context, adScene *model.AdScene) error {
	err := r.tx.DB(ctx).Model(adScene).Updates(adScene).Error
	if err != nil {
		return err
	}

	return err
}

// Delete 批量删除
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	var adScenes []model.AdScene
	for _, id := range ids {
		// 根据ID获取标签
		adScene, err := r.Detail(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的推广场景", id)
		}
		adScenes = append(adScenes, adScene)
	}

	err := r.tx.DB(ctx).Unscoped().Delete(&adScenes).Error

	return err
}
