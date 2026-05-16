// Package repository implements tag lookup and post-tag relation query logic.
package repository

// 本文件实现标签查询和帖子标签关联查询能力。

import (
	"context"
	"strings"

	"gotribe/internal/core/database"
	tagmodel "gotribe/internal/model"
)

// Repository 负责封装标签相关的数据访问逻辑。
type Repository struct {
	tx *database.TransactionManager
}

// List 获取所有启用的标签，按 sort desc 排序。
func (r *Repository) List(ctx context.Context, keyword string, limit int) ([]tagmodel.Tag, error) {
	db := r.tx.DB(ctx).Model(&tagmodel.Tag{}).Where("status = 1")
	if keyword != "" {
		db = db.Where("title LIKE ? OR slug LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if limit > 0 {
		db = db.Limit(limit)
	}
	var tags []tagmodel.Tag
	if err := db.Order("sort DESC, id ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// postTagRow 表示文章与标签关联表的查询结果。
type postTagRow struct {
	PostID int64 `gorm:"column:post_id"`
	TagID  int64 `gorm:"column:tag_id"`
}

// NewRepository 创建标签仓储实例。
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// FindIDsByKeyword 按标签 slug 或标题查找匹配的标签 ID。
func (r *Repository) FindIDsByKeyword(ctx context.Context, keyword string) ([]int64, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, nil
	}

	var tags []tagmodel.Tag
	if err := r.tx.DB(ctx).
		Model(&tagmodel.Tag{}).
		Where("(slug = ? OR title = ?)", keyword, keyword).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, nil
	}

	tagIDs := make([]int64, 0, len(tags))
	for _, item := range tags {
		tagIDs = append(tagIDs, item.ID)
	}
	return tagIDs, nil
}

// ListByPostIDs 按文章 ID 集合批量加载标签列表。
func (r *Repository) ListByPostIDs(ctx context.Context, postIDs []int64) (map[int64][]tagmodel.Tag, error) {
	result := make(map[int64][]tagmodel.Tag, len(postIDs))
	if len(postIDs) == 0 {
		return result, nil
	}

	var relations []postTagRow
	if err := r.tx.DB(ctx).
		Table("post_tag").
		Where("post_id IN ?", postIDs).
		Order("post_id asc").
		Order("tag_id asc").
		Find(&relations).Error; err != nil {
		return nil, err
	}
	if len(relations) == 0 {
		return result, nil
	}

	tagIDs := make([]int64, 0, len(relations))
	seenTagIDs := make(map[int64]struct{}, len(relations))
	for _, relation := range relations {
		if _, exists := seenTagIDs[relation.TagID]; exists {
			continue
		}
		seenTagIDs[relation.TagID] = struct{}{}
		tagIDs = append(tagIDs, relation.TagID)
	}

	var tags []tagmodel.Tag
	if err := r.tx.DB(ctx).
		Model(&tagmodel.Tag{}).
		Where("id IN ?", tagIDs).
		Order("sort desc").
		Order("id asc").
		Find(&tags).Error; err != nil {
		return nil, err
	}

	tagsByID := make(map[int64]tagmodel.Tag, len(tags))
	for _, tag := range tags {
		tagsByID[tag.ID] = tag
	}

	for _, relation := range relations {
		tag, exists := tagsByID[relation.TagID]
		if !exists {
			continue
		}
		result[relation.PostID] = append(result[relation.PostID], tag)
	}
	return result, nil
}
