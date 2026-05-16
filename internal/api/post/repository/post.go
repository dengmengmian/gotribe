package repository

// 本文件实现帖子查询相关的数据访问逻辑。

import (
	"context"
	"strings"

	"gotribe/internal/core/database"
	postmodel "gotribe/internal/model"

	"gorm.io/gorm"
)

// Repository 负责封装文章相关的数据访问逻辑。
type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建文章仓储实例。
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// List 按条件查询文章列表并返回总数。
func (r *Repository) List(ctx context.Context, projectID string, filter ListFilter) ([]postmodel.Post, int64, error) {
	var (
		posts []postmodel.Post
		total int64
	)

	normalizedPage, normalizedPerPage := database.NormalizePagination(filter.Page, filter.PerPage)
	db := r.baseQuery(ctx, projectID, filter)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("is_top desc").
		Order("sort desc").
		Order("created_at desc").
		Offset((normalizedPage - 1) * normalizedPerPage).
		Limit(normalizedPerPage).
		Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

// GetByPostID 按业务文章 ID 或 URL 别名读取文章详情。
func (r *Repository) GetByPostID(ctx context.Context, projectID, postID string) (*postmodel.Post, error) {
	var post postmodel.Post
	db := r.tx.DB(ctx).Model(&postmodel.Post{}).Where("(post_id = ? OR slug = ?)", postID, postID)
	if projectID != "" {
		db = db.Where("project_id = ?", projectID)
	}
	if err := db.First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// ListByPostIDs 按业务文章 ID 集合批量读取文章。
func (r *Repository) ListByPostIDs(ctx context.Context, projectID string, postIDs []string) ([]postmodel.Post, error) {
	posts := make([]postmodel.Post, 0, len(postIDs))
	if len(postIDs) == 0 {
		return posts, nil
	}

	db := r.tx.DB(ctx).Model(&postmodel.Post{}).Where("post_id IN ?", postIDs)
	if projectID != "" {
		db = db.Where("project_id = ?", projectID)
	}
	if err := db.Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// baseQuery 构建文章列表查询的基础条件。
func (r *Repository) baseQuery(ctx context.Context, projectID string, filter ListFilter) *gorm.DB {
	db := r.tx.DB(ctx).Model(&postmodel.Post{})
	if projectID != "" {
		db = db.Where("project_id = ?", projectID)
	}
	if len(filter.TagIDs) > 0 {
		subQuery := r.tx.DB(ctx).
			Table("post_tag").
			Select("distinct post_id").
			Where("tag_id IN ?", filter.TagIDs)
		db = db.Where("id IN (?)", subQuery)
	}
	if filter.Status != nil {
		db = db.Where("status = ?", *filter.Status)
	} else {
		db = db.Where("status = ?", 2)
	}
	if filter.Type != nil {
		db = db.Where("type = ?", *filter.Type)
	}
	if strings.TrimSpace(filter.Keyword) != "" {
		like := "%" + strings.TrimSpace(filter.Keyword) + "%"
		db = db.Where("(title ILIKE ? OR description ILIKE ?)", like, like)
	}
	if strings.TrimSpace(filter.DynamicType) != "" {
		db = db.Where("dynamic_type = ?", strings.TrimSpace(filter.DynamicType))
	}
	return db
}
