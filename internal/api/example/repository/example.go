// Package repository implements data access logic for example module operations.
package repository

// 本文件实现 example 模块的数据访问逻辑。

import (
	"context"
	"strings"

	"gotribe/internal/core/database"
	examplemodel "gotribe/internal/model"

	"gorm.io/gorm"
)

// Repository 负责封装示例业务单相关的数据访问逻辑。
type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建示例业务单仓储实例。
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// Create 创建示例业务单主记录。
func (r *Repository) Create(ctx context.Context, entity *examplemodel.Example) error {
	return r.tx.DB(ctx).Create(entity).Error
}

// GetByExampleID 按业务 ID 读取示例业务单。
func (r *Repository) GetByExampleID(ctx context.Context, projectID string, userID int64, exampleID string) (*examplemodel.Example, error) {
	var entity examplemodel.Example
	if err := r.scopedExamples(ctx, projectID, userID).
		Where("example_id = ?", exampleID).
		First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

// List 按条件查询示例业务单列表并返回总数。
func (r *Repository) List(ctx context.Context, projectID string, userID int64, filter ListFilter) ([]examplemodel.Example, int64, error) {
	var (
		items []examplemodel.Example
		total int64
	)

	page, perPage := database.NormalizePagination(filter.Page, filter.PerPage)
	db := r.baseListQuery(ctx, projectID, userID, filter)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("created_at desc").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// UpdateByID 按主键更新示例业务单。
func (r *Repository) UpdateByID(ctx context.Context, id int64, updates map[string]any) error {
	return r.tx.DB(ctx).
		Model(&examplemodel.Example{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// DeleteByID 按主键删除示例业务单。
func (r *Repository) DeleteByID(ctx context.Context, id int64) error {
	return r.tx.DB(ctx).Delete(&examplemodel.Example{}, id).Error
}

// CreatePosts 批量写入示例业务单关联文章。
func (r *Repository) CreatePosts(ctx context.Context, items []examplemodel.ExamplePost) error {
	if len(items) == 0 {
		return nil
	}
	return r.tx.DB(ctx).Create(&items).Error
}

// DeletePostsByExampleRecordID 删除指定示例业务单的全部关联文章。
func (r *Repository) DeletePostsByExampleRecordID(ctx context.Context, exampleRecordID int64) error {
	return r.tx.DB(ctx).
		Where("example_record_id = ?", exampleRecordID).
		Delete(&examplemodel.ExamplePost{}).Error
}

// ListPostsByExampleRecordIDs 批量读取示例业务单的关联文章。
func (r *Repository) ListPostsByExampleRecordIDs(ctx context.Context, exampleRecordIDs []int64) (map[int64][]examplemodel.ExamplePost, error) {
	result := make(map[int64][]examplemodel.ExamplePost, len(exampleRecordIDs))
	if len(exampleRecordIDs) == 0 {
		return result, nil
	}

	var rows []examplemodel.ExamplePost
	if err := r.tx.DB(ctx).
		Model(&examplemodel.ExamplePost{}).
		Where("example_record_id IN ?", exampleRecordIDs).
		Order("example_record_id asc").
		Order("sort asc").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.ExampleRecordID] = append(result[row.ExampleRecordID], row)
	}
	return result, nil
}

func (r *Repository) baseListQuery(ctx context.Context, projectID string, userID int64, filter ListFilter) *gorm.DB {
	db := r.scopedExamples(ctx, projectID, userID)
	if filter.Status != nil {
		db = db.Where("status = ?", *filter.Status)
	}
	if strings.TrimSpace(filter.Keyword) != "" {
		like := "%" + strings.TrimSpace(filter.Keyword) + "%"
		db = db.Where("(name ILIKE ? OR description ILIKE ?)", like, like)
	}
	return db
}

func (r *Repository) scopedExamples(ctx context.Context, projectID string, userID int64) *gorm.DB {
	return r.tx.DB(ctx).
		Model(&examplemodel.Example{}).
		Where("project_id = ? AND user_id = ?", projectID, userID)
}
