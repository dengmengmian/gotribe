package repository

import (
	"context"

	"gotribe/internal/admin/feedback/dto"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewFeedbackRepository 创建反馈仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// List 获取反馈列表
func (r *Repository) List(ctx context.Context, req *dto.FeedbackListRequest) ([]*model.Feedback, int64, error) {
	var list []*model.Feedback
	db := r.tx.DB(ctx).Model(&model.Feedback{})

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

	return list, total, err
}

// Delete 批量删除反馈。
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	return r.tx.DB(ctx).Where("id IN ?", ids).Delete(&model.Feedback{}).Error
}
