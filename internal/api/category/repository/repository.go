// Package repository implements category query logic for the public API.
package repository

import (
	"context"

	"gotribe/internal/core/database"
	"gotribe/internal/model"
)

// Repository 负责封装分类相关的数据访问逻辑。
type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建分类仓储实例。
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// ListByParentID 按父分类 ID 查询所有启用的子分类，按 sort 排序。
func (r *Repository) ListByParentID(ctx context.Context, parentID int64) ([]model.Category, error) {
	var categories []model.Category
	err := r.tx.DB(ctx).
		Model(&model.Category{}).
		Where("parent_id = ? AND status = 1 AND hidden = 1", parentID).
		Order("sort ASC, id ASC").
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
