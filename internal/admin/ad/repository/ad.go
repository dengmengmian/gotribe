package repository

import (
	"context"
	"fmt"
	"strings"

	"gotribe/internal/admin/ad/dto"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
	"gotribe/internal/core/util"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewAdRepository 创建广告仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

func buildAdOrder(req *dto.AdListRequest) string {
	sortByMap := map[string]string{
		"id":          "id",
		"title":       "title",
		"description": "description",
		"status":      "status",
		"createdAt":   "created_at",
		"created_at":  "created_at",
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

// Detail 获取单个推广场景
func (r *Repository) Detail(ctx context.Context, id int64) (model.Ad, error) {
	var ad model.Ad
	err := r.tx.DB(ctx).Where("id = ?", id).First(&ad).Error
	return ad, err
}

// List 获取推广场景列表
func (r *Repository) List(ctx context.Context, req *dto.AdListRequest) ([]*model.Ad, int64, error) {
	var list []*model.Ad
	db := r.tx.DB(ctx).Model(&model.Ad{})

	if req.SceneID > 0 {
		db = db.Where("scene_id = ?", req.SceneID)
	}
	if !utils.IsEmpty(req.Title) {
		db = db.Where("title like ?", fmt.Sprintf("%%%s%%", req.Title))
	}
	if !utils.IsEmpty(req.Status) {
		db = db.Where("status = ?", fmt.Sprintf("%d", req.Status))
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	db = db.Order(buildAdOrder(req))

	page, perPage := database.NormalizePagination(int(req.PageNum), int(req.PageSize))
	err = db.Offset((page - 1) * perPage).Limit(perPage).Find(&list).Error

	list, err = r.GetAdOther(ctx, list)
	return list, total, err
}

// GetAdOther 获取推广场景其他信息
func (r *Repository) GetAdOther(ctx context.Context, ads []*model.Ad) ([]*model.Ad, error) {
	if len(ads) == 0 {
		return ads, nil
	}

	// 收集所有 SceneID
	sceneIDs := make([]int64, 0, len(ads))
	for _, m := range ads {
		if m.SceneID > 0 {
			sceneIDs = append(sceneIDs, m.SceneID)
		}
	}

	if len(sceneIDs) == 0 {
		return ads, nil
	}

	// 批量查询
	var adScenes []*model.AdScene
	if err := r.tx.DB(ctx).Where("id IN ?", sceneIDs).Find(&adScenes).Error; err != nil {
		return ads, err
	}

	// 建立映射
	adSceneMap := make(map[int64]*model.AdScene)
	for _, scene := range adScenes {
		adSceneMap[scene.ID] = scene
	}

	// 赋值
	for _, m := range ads {
		if scene, ok := adSceneMap[m.SceneID]; ok {
			m.Scene = scene
		}
	}
	return ads, nil
}

// Create 创建推广场景
func (r *Repository) Create(ctx context.Context, ad *model.Ad) error {
	err := r.tx.DB(ctx).Create(ad).Error
	return err
}

// Update 更新推广场景
func (r *Repository) Update(ctx context.Context, ad *model.Ad) error {
	err := r.tx.DB(ctx).Model(ad).Updates(ad).Error
	if err != nil {
		return err
	}

	return err
}

// Delete 批量删除
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	var ads []model.Ad
	for _, id := range ids {
		// 根据ID获取标签
		ad, err := r.Detail(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的推广场景", id)
		}
		ads = append(ads, ad)
	}

	err := r.tx.DB(ctx).Unscoped().Delete(&ads).Error

	return err
}
