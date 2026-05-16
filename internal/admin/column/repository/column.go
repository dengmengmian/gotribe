package repository

import (
	"context"
	"fmt"
	"strings"

	"gotribe/internal/admin/column/dto"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建专栏仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// Detail 获取单个专栏
func (r *Repository) Detail(ctx context.Context, id int64) (model.Column, error) {
	var column model.Column
	err := r.tx.DB(ctx).Where("id = ?", id).First(&column).Error
	return column, err
}

// List 获取专栏列表
func (r *Repository) List(ctx context.Context, req *dto.ColumnListRequest) ([]*model.Column, int64, error) {
	var list []*model.Column
	db := r.tx.DB(ctx).Model(&model.Column{})

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
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

	return list, total, err
}

// Create 创建专栏
func (r *Repository) Create(ctx context.Context, column *model.Column) error {
	err := r.tx.DB(ctx).Create(column).Error
	return err
}

// Update 更新专栏
func (r *Repository) Update(ctx context.Context, column *model.Column) error {
	err := r.tx.DB(ctx).Model(column).Updates(column).Error
	if err != nil {
		return err
	}

	return err
}

// Delete 批量删除
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	var columns []model.Column
	for _, id := range ids {
		// 根据ID获取专栏
		column, err := r.Detail(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的专栏", id)
		}
		columns = append(columns, column)
	}

	err := r.tx.DB(ctx).Delete(&columns).Error

	return err
}
