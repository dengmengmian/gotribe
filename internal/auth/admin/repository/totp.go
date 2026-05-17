// Package repository provides data access for admin authentication features.
package repository

// 本文件实现 admin_totp 表的 CRUD。所有 secret 的加解密由 service 层负责，
// repository 只做存储字段透传。

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	"gotribe/internal/model"
)

// TOTPRepository TOTP 绑定记录数据访问。
type TOTPRepository struct {
	tx *database.TransactionManager
}

// NewTOTPRepository 创建仓库实例。
func NewTOTPRepository(tx *database.TransactionManager) *TOTPRepository {
	return &TOTPRepository{tx: tx}
}

// GetByAdminID 根据 admin_id 查询 TOTP 记录（含未启用）。无记录返回 nil, nil。
func (r *TOTPRepository) GetByAdminID(ctx context.Context, adminID int64) (*model.AdminTOTP, error) {
	var record model.AdminTOTP
	err := r.tx.DB(ctx).Where("admin_id = ?", adminID).First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.Internal("查询 TOTP 记录失败", err)
	}
	return &record, nil
}

// Upsert 创建或更新一条记录（按 admin_id 唯一）。
// 用于 Bind 阶段（写入未启用的 secret）。
func (r *TOTPRepository) Upsert(ctx context.Context, record *model.AdminTOTP) error {
	existing, err := r.GetByAdminID(ctx, record.AdminID)
	if err != nil {
		return err
	}
	if existing == nil {
		if err := r.tx.DB(ctx).Create(record).Error; err != nil {
			return errs.Internal("创建 TOTP 记录失败", err)
		}
		return nil
	}
	record.ID = existing.ID
	record.CreatedAt = existing.CreatedAt
	if err := r.tx.DB(ctx).Save(record).Error; err != nil {
		return errs.Internal("更新 TOTP 记录失败", err)
	}
	return nil
}

// MarkEnabled 将记录置为启用并更新 secret / recovery_codes（如有变更则一并写）。
func (r *TOTPRepository) MarkEnabled(ctx context.Context, adminID int64) error {
	res := r.tx.DB(ctx).Model(&model.AdminTOTP{}).
		Where("admin_id = ?", adminID).
		Update("enabled", true)
	if res.Error != nil {
		return errs.Internal("启用 TOTP 失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return errs.NotFound("TOTP 记录不存在", nil)
	}
	return nil
}

// UpdateRecoveryCodes 替换备份码 JSON。
func (r *TOTPRepository) UpdateRecoveryCodes(ctx context.Context, adminID int64, codesJSON string) error {
	res := r.tx.DB(ctx).Model(&model.AdminTOTP{}).
		Where("admin_id = ?", adminID).
		Update("recovery_codes", codesJSON)
	if res.Error != nil {
		return errs.Internal("更新备份码失败", res.Error)
	}
	return nil
}

// UpdateLastUsedAt 记录最近成功使用 TOTP 的时间戳（秒）。
func (r *TOTPRepository) UpdateLastUsedAt(ctx context.Context, adminID, unixSeconds int64) error {
	return r.tx.DB(ctx).Model(&model.AdminTOTP{}).
		Where("admin_id = ?", adminID).
		Update("last_used_at", unixSeconds).Error
}

// DeleteByAdminID 物理删除 TOTP 记录（用于自助解绑或超管重置）。
func (r *TOTPRepository) DeleteByAdminID(ctx context.Context, adminID int64) error {
	return r.tx.DB(ctx).Unscoped().Where("admin_id = ?", adminID).Delete(&model.AdminTOTP{}).Error
}
