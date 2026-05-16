// Package repository implements data access logic for user event storage.
package repository

// 本文件实现用户行为写入的数据访问逻辑。

import (
	"context"

	"gotribe/internal/core/database"
	usereventmodel "gotribe/internal/model"
)

// Repository 负责封装用户行为事件相关的数据访问逻辑。
type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建用户行为事件仓储实例。
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// Create 写入用户行为事件。
func (r *Repository) Create(ctx context.Context, event *usereventmodel.UserEvent) error {
	return r.tx.DB(ctx).Create(event).Error
}
