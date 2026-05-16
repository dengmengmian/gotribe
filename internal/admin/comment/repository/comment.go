package repository

import (
	"context"
	"fmt"
	"strings"

	"gotribe/internal/admin/comment/dto"
	"gotribe/internal/core/database"
	"gotribe/internal/core/util"
	"gotribe/internal/model"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建评论仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

func buildCommentOrder(req *dto.CommentListRequest) string {
	sortByMap := map[string]string{
		"id":         "id",
		"comment":    "content",
		"userID":     "user_id",
		"user_id":    "user_id",
		"status":     "status",
		"createdAt":  "created_at",
		"created_at": "created_at",
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

// Detail 获取单个评论
func (r *Repository) Detail(ctx context.Context, id int64) (model.Comment, error) {
	var comment model.Comment
	err := r.tx.DB(ctx).Where("id = ?", id).First(&comment).Error
	return comment, err
}

// List 获取评论列表
func (r *Repository) List(ctx context.Context, req *dto.CommentListRequest) ([]*model.Comment, int64, error) {
	var list []*model.Comment
	db := r.tx.DB(ctx).Model(&model.Comment{})

	objectID := strings.TrimSpace(req.ObjectID)
	if !utils.IsEmpty(objectID) {
		db = db.Where("object_id = ?", objectID)
	}
	if !utils.IsEmpty(req.ObjectType) {
		db = db.Where("object_type = ?", req.ObjectType)
	}
	if !utils.IsEmpty(req.Status) {
		db = db.Where("status = ?", req.Status)
	}
	if req.ProjectID > 0 {
		db = db.Where("project_id = ?", req.ProjectID)
	}
	if !utils.IsEmpty(req.Nickname) {
		nicknameIDs := r.tx.DB(ctx).
			Model(&model.User{}).
			Select("id").
			Where("nickname like ?", fmt.Sprintf("%%%s%%", req.Nickname))
		db = db.Where("user_id IN (?)", nicknameIDs)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	db = db.Order(buildCommentOrder(req))

	page, perPage := database.NormalizePagination(int(req.PageNum), int(req.PageSize))
	err = db.Offset((page - 1) * perPage).Limit(perPage).Find(&list).Error

	list, err = r.GetCommentOther(ctx, list)
	return list, total, err
}

// GetCommentOther 获取评论其他信息
func (r *Repository) GetCommentOther(ctx context.Context, comments []*model.Comment) ([]*model.Comment, error) {
	if len(comments) == 0 {
		return comments, nil
	}

	// 收集所有 UserID
	userIDs := make([]int64, 0, len(comments))
	for _, m := range comments {
		if m.UserID > 0 {
			userIDs = append(userIDs, m.UserID)
		}
	}

	if len(userIDs) == 0 {
		return comments, nil
	}

	// 批量查询
	var users []*model.User
	if err := r.tx.DB(ctx).Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return comments, err
	}

	// 建立映射
	userMap := make(map[int64]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 赋值
	for _, m := range comments {
		if user, ok := userMap[m.UserID]; ok {
			m.User = user
		}
	}
	return comments, nil
}

// Delete 删除评论。
func (r *Repository) Delete(ctx context.Context, id int64) error {
	return r.tx.DB(ctx).Where("id = ?", id).Delete(&model.Comment{}).Error
}

// Update 更新评论。
func (r *Repository) Update(ctx context.Context, comment *model.Comment) error {
	err := r.tx.DB(ctx).Model(comment).Updates(comment).Error
	if err != nil {
		return err
	}
	return err
}
