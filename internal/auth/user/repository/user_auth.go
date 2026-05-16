// Package repository implements data access logic for authentication-related queries.
package repository

// 本文件实现认证模块查询用户身份信息的数据访问逻辑。

import (
	"context"

	authmodel "gotribe/internal/model"
	"gotribe/internal/core/database"
)

// UserAuthRepository 负责封装认证相关的数据访问逻辑。
type UserAuthRepository struct {
	tx *database.TransactionManager
}

// NewUserAuthRepository 创建认证模块使用的用户仓储。
func NewUserAuthRepository(tx *database.TransactionManager) *UserAuthRepository {
	return &UserAuthRepository{tx: tx}
}

// FindByIdentity 按用户名、邮箱或手机号查找可登录用户。
func (r *UserAuthRepository) FindByIdentity(ctx context.Context, projectID, identity string) (*authmodel.AuthUser, error) {
	var user authmodel.AuthUser
	db := r.tx.DB(ctx).Model(&authmodel.AuthUser{}).
		Where("(username = ? OR email = ? OR phone = ?)", identity, identity, identity)
	if projectID != "" {
		db = db.Where("project_id = ?", projectID)
	}
	if err := db.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 按用户 ID 读取认证所需的用户信息。
func (r *UserAuthRepository) FindByID(ctx context.Context, projectID string, userID int64) (*authmodel.AuthUser, error) {
	var user authmodel.AuthUser
	db := r.tx.DB(ctx).Model(&authmodel.AuthUser{}).Where("id = ?", userID)
	if projectID != "" {
		db = db.Where("project_id = ?", projectID)
	}
	if err := db.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
