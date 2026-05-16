package repository

import (
	"context"
	"errors"

	"gotribe/internal/model"

	"gotribe/internal/core/database"

	"gorm.io/gorm"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewSystemConfigRepository 创建系统配置仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// Detail 获取系统配置
func (r *Repository) Detail(ctx context.Context) (model.SystemConfig, error) {
	var systemConfig model.SystemConfig
	err := r.tx.DB(ctx).First(&systemConfig).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return systemConfig, nil
	}
	return systemConfig, err
}

// Update 更新系统配置
func (r *Repository) Update(ctx context.Context, systemConfig *model.SystemConfig) error {
	if systemConfig.ID == 0 {
		return r.tx.DB(ctx).Create(systemConfig).Error
	}
	err := r.tx.DB(ctx).Model(systemConfig).Updates(systemConfig).Error
	if err != nil {
		return err
	}

	return err
}
