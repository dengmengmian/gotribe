package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gotribe/internal/admin/tag/dto"
	"gotribe/internal/model"

	"gorm.io/gorm"
	"gotribe/internal/core/database"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建标签仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

func buildTagOrder(req *dto.TagListRequest) string {
	sortByMap := map[string]string{
		"id":          "id",
		"title":       "title",
		"description": "description",
		"createdAt":   "created_at",
		"created_at":  "created_at",
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

// Detail 获取单个标签
func (r *Repository) Detail(ctx context.Context, id int64) (model.Tag, error) {
	var tag model.Tag
	err := r.tx.DB(ctx).Where("id = ?", id).First(&tag).Error
	return tag, err
}

// List 获取标签列表
func (r *Repository) List(ctx context.Context, req *dto.TagListRequest) ([]*model.Tag, int64, error) {
	var list []*model.Tag
	db := r.tx.DB(ctx).Model(&model.Tag{})

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	id := req.ID
	if id > 0 {
		db = db.Where("id = ?", id)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	db = db.Order(buildTagOrder(req))

	page, perPage := database.NormalizePagination(int(req.PageNum), int(req.PageSize))
	err = db.Offset((page - 1) * perPage).Limit(perPage).Find(&list).Error

	return list, total, err
}

// Create 创建标签
func (r *Repository) Create(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	if r.isTagExist(ctx, tag.Title) {
		return nil, fmt.Errorf("%s标签已存在", tag.Title)
	}
	result := r.tx.DB(ctx).Create(tag)
	if result.Error != nil {
		return nil, result.Error
	}
	return tag, nil
}

// Update 更新标签
func (r *Repository) Update(ctx context.Context, tag *model.Tag) error {
	err := r.tx.DB(ctx).Model(tag).Updates(tag).Error
	if err != nil {
		return err
	}
	return err
}

// Delete 批量删除标签
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	var tags []model.Tag
	for _, id := range ids {
		// 根据ID获取标签
		tag, err := r.Detail(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的标签", id)
		}
		tags = append(tags, tag)
	}

	err := r.tx.DB(ctx).Unscoped().Delete(&tags).Error
	return err
}

func (r *Repository) isTagExist(ctx context.Context, title string) bool {
	var tag model.Tag
	result := r.tx.DB(ctx).Where("title = ?", title).First(&tag)
	return !errors.Is(result.Error, gorm.ErrRecordNotFound)
}
