// Package repository implements data access logic for user profile operations.
package repository

// 本文件实现当前用户资料相关的数据访问逻辑。

import (
	"context"

	"gotribe/internal/core/database"
	profilemodel "gotribe/internal/model"
)

// Repository 负责封装当前用户资料相关的数据访问逻辑。
type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建当前用户资料仓储实例。
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// GetByID 根据用户 ID 读取用户资料。
func (r *Repository) GetByID(ctx context.Context, projectID string, userID int64) (*profilemodel.UserProfile, error) {
	var user profilemodel.UserProfile
	db := r.tx.DB(ctx).Model(&profilemodel.UserProfile{}).Where("id = ?", userID)
	if projectID != "" {
		db = db.Where("project_id = ?", projectID)
	}
	if err := db.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetPasswordByID 根据用户 ID 读取当前密码哈希。
func (r *Repository) GetPasswordByID(ctx context.Context, projectID string, userID int64) (string, error) {
	var user profilemodel.UserProfile
	db := r.tx.DB(ctx).Select("id", "password").Model(&profilemodel.UserProfile{}).Where("id = ?", userID)
	if projectID != "" {
		db = db.Where("project_id = ?", projectID)
	}
	if err := db.First(&user).Error; err != nil {
		return "", err
	}
	return user.Password, nil
}

// Update 更新用户资料字段。
func (r *Repository) Update(ctx context.Context, projectID string, userID int64, updates map[string]any) error {
	db := r.tx.DB(ctx).Model(&profilemodel.UserProfile{}).Where("id = ?", userID)
	if projectID != "" {
		db = db.Where("project_id = ?", projectID)
	}
	return db.Updates(updates).Error
}

// UpdatePassword 更新用户密码哈希。
func (r *Repository) UpdatePassword(ctx context.Context, projectID string, userID int64, hashedPassword string) error {
	db := r.tx.DB(ctx).Model(&profilemodel.UserProfile{}).Where("id = ?", userID)
	if projectID != "" {
		db = db.Where("project_id = ?", projectID)
	}
	return db.Update("password", hashedPassword).Error
}
