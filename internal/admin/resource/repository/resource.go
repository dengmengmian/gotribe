package repository

import (
	"context"
	"fmt"

	"gotribe/internal/admin/resource/dto"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建资源仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// Detail 获取单个资源
func (r *Repository) Detail(ctx context.Context, id int64) (model.Resource, error) {
	var resource model.Resource
	err := r.tx.DB(ctx).Where("id = ?", id).First(&resource).Error
	return resource, err
}

// List 获取资源列表
func (r *Repository) List(ctx context.Context, req *dto.ResourceListRequest) ([]*model.Resource, int64, error) {
	var list []*model.Resource
	db := r.tx.DB(ctx).Model(&model.Resource{})

	if int(req.Type) > 0 {
		db = db.Where("file_type = ?", req.Type)
	}

	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
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

	return list, total, err
}

// Create 创建资源
func (r *Repository) Create(ctx context.Context, resource *model.Resource) error {
	err := r.tx.DB(ctx).Create(resource).Error
	return err
}

// Update 更新资源
func (r *Repository) Update(ctx context.Context, resource *model.Resource) error {
	err := r.tx.DB(ctx).Model(resource).Updates(resource).Error
	if err != nil {
		return err
	}

	return err
}

// Delete 删除文件
func (r *Repository) Delete(ctx context.Context, id int64) error {
	resource, err := r.Detail(ctx, id)
	if err != nil {
		return fmt.Errorf("未获取到ID为%d的资源", id)
	}

	// 硬删除（CDN 删除由 service 层编排）
	return r.tx.DB(ctx).Unscoped().Delete(&resource).Error
}
