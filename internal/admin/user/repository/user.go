package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gotribe/internal/admin/user/dto"
	"gotribe/internal/model"
	"strings"

	"gorm.io/gorm"
	"gotribe/internal/core/database"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建用户仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// Detail 获取单个用户
func (r *Repository) Detail(ctx context.Context, id int64) (model.User, error) {
	var user model.User
	err := r.tx.DB(ctx).Where("id = ?", id).First(&user).Error
	return user, err
}

// UserProfileRef 承载失效 ToC 用户资料缓存所需的最小字段。
type UserProfileRef struct {
	ID        int64
	ProjectID int64
}

// ListProfileRefsByIDs 按用户 ID 列表查询失效 ToC profile 缓存所需的 id / project_id。
func (r *Repository) ListProfileRefsByIDs(ctx context.Context, ids []int64) ([]UserProfileRef, error) {
	refs := make([]UserProfileRef, 0, len(ids))
	if len(ids) == 0 {
		return refs, nil
	}
	err := r.tx.DB(ctx).Model(&model.User{}).
		Select("id", "project_id").
		Where("id IN ?", ids).
		Find(&refs).Error
	return refs, err
}

// List 获取用户列表
func (r *Repository) List(ctx context.Context, req *dto.UserListRequest) ([]*model.User, int64, error) {
	var list []*model.User
	db := r.tx.DB(ctx).Model(&model.User{})

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	if req.UserID > 0 {
		db = db.Where("id = ?", req.UserID)
	}
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

	return r.GetUserOther(ctx, list), total, err
}

func (r *Repository) GetUserOther(ctx context.Context, user []*model.User) []*model.User {
	for _, m := range user {
		userPoint := r.GetUserPoint(ctx, m.ID)
		m.Point = userPoint
	}
	return user
}

func (r *Repository) GetUserPoint(ctx context.Context, userID int64) float64 {
	var sum sql.NullFloat64
	var pointAvailable *model.PointAvailable
	row := r.tx.DB(ctx).Model(&pointAvailable).Select("SUM(points)").Where("user_id = ?", userID).Row()
	err := row.Scan(&sum)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果记录不存在，返回 0
			return 0
		}
		return 0
	}
	// 如果 sum 是 NULL，则返回 0
	if !sum.Valid {
		return 0
	}
	return sum.Float64
}

// Create 创建用户
func (r *Repository) Create(ctx context.Context, user *model.User) error {
	normalizeUserContactFields(user)
	err := r.tx.DB(ctx).Create(user).Error
	return err
}

// Update 更新用户
func (r *Repository) Update(ctx context.Context, user *model.User) error {
	normalizeUserContactFields(user)
	err := r.tx.DB(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"username":   user.Username,
		"project_id": user.ProjectID,
		"password":   user.Password,
		"nickname":   user.Nickname,
		"email":      user.Email,
		"phone":      user.Phone,
		"sex":        user.Sex,
		"status":     user.Status,
		"birthday":   user.Birthday,
		"background": user.Background,
		"ext":        user.Ext,
		"avatar_url": user.AvatarURL,
	}).Error
	if err != nil {
		return err
	}

	return err
}

// Delete 批量删除
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	var users []model.User
	for _, id := range ids {
		// 根据ID获取用户
		user, err := r.Detail(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的用户", id)
		}
		users = append(users, user)
	}

	err := r.tx.DB(ctx).Delete(&users).Error

	return err
}

// Search 根据昵称搜索用户
func (r *Repository) Search(ctx context.Context, nickname string) ([]*model.User, error) {
	var list []*model.User
	db := r.tx.DB(ctx).Model(&model.User{}).Order("created_at DESC")

	if strings.TrimSpace(nickname) != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	err := db.Find(&list).Error
	return list, err
}

func normalizeUserContactFields(user *model.User) {
	if user == nil {
		return
	}
	user.Email = trimOptionalString(user.Email)
	user.Phone = trimOptionalString(user.Phone)
}

func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
