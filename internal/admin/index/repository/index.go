package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gotribe/internal/admin/index/dto"
	"gotribe/internal/core/database"
)

const defaultProjectID = "1"

// Repository 仪表盘数据访问。
type Repository struct {
	tx *database.TransactionManager
}

// NewRepository 创建仪表盘数据访问实例。
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// Stats 返回统计卡片数据。
func (r *Repository) Stats(ctx context.Context, projectID string) (dto.Stats, error) {
	projectID = normalizeProjectID(projectID)
	db := r.tx.DB(ctx)
	weekAgo := time.Now().AddDate(0, 0, -7)

	var total int64
	if err := db.Table("post").
		Where("project_id = ? AND deleted_at IS NULL", projectID).
		Count(&total).Error; err != nil {
		return dto.Stats{}, err
	}

	var drafts int64
	if err := db.Table("post").
		Where("project_id = ? AND status = 1 AND deleted_at IS NULL", projectID).
		Count(&drafts).Error; err != nil {
		return dto.Stats{}, err
	}

	var pending int64
	if err := db.Table("comment").
		Where("project_id = ? AND status = 1 AND deleted_at IS NULL", projectID).
		Count(&pending).Error; err != nil {
		return dto.Stats{}, err
	}

	var visits int64
	if err := db.Table("user_event").
		Where("event_type = 1 AND created_at >= ? AND project_id = ?", weekAgo, projectID).
		Count(&visits).Error; err != nil {
		return dto.Stats{}, err
	}

	return dto.Stats{
		TotalPosts:      total,
		DraftPosts:      drafts,
		PendingComments: pending,
		WeekVisits:      visits,
	}, nil
}

// VisitTrend 近 7 日每日访问量。
func (r *Repository) VisitTrend(ctx context.Context, projectID string) ([]dto.VisitPoint, error) {
	projectID = normalizeProjectID(projectID)
	db := r.tx.DB(ctx)
	startDate := time.Now().AddDate(0, 0, -6)

	type row struct {
		Date  string `gorm:"column:date"`
		Count int64  `gorm:"column:count"`
	}

	var rows []row
	if err := db.Table("user_event").
		Select("TO_CHAR(created_at, 'MM-DD') as date, COUNT(*) as count").
		Where("event_type = 1 AND created_at >= ? AND project_id = ?", startDate, projectID).
		Group("TO_CHAR(created_at, 'MM-DD')").
		Order("date").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	points := make([]dto.VisitPoint, 0, len(rows))
	for _, r := range rows {
		points = append(points, dto.VisitPoint{
			Date:      r.Date,
			Visits:    r.Count,
			PageViews: r.Count * 2,
		})
	}
	return points, nil
}

// RecentPosts 最近 N 篇文章。
func (r *Repository) RecentPosts(ctx context.Context, projectID string, limit int) ([]dto.PostSummary, error) {
	projectID = normalizeProjectID(projectID)
	db := r.tx.DB(ctx)

	var rows []dto.PostSummary
	if err := db.Table("post").
		Select("id, title, status, created_at, view").
		Where("project_id = ? AND deleted_at IS NULL", projectID).
		Order("created_at DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// RecentComments 最近 N 条评论。
func (r *Repository) RecentComments(ctx context.Context, projectID string, limit int) ([]dto.CommentSummary, error) {
	projectID = normalizeProjectID(projectID)
	db := r.tx.DB(ctx)

	var rows []struct {
		ID        int64     `gorm:"column:id"`
		Nickname  string    `gorm:"column:nickname"`
		Content   string    `gorm:"column:content"`
		CreatedAt time.Time `gorm:"column:created_at"`
		PostTitle string    `gorm:"column:post_title"`
	}

	if err := db.Table("comment").
		Select(`comment.id, COALESCE("user".nickname, '匿名') as nickname, comment.content, comment.created_at, COALESCE(post.title, '') as post_title`).
		Joins(`LEFT JOIN "user" ON "user".id = comment.user_id`).
		Joins("LEFT JOIN post ON post.post_id = comment.object_id AND comment.object_type = 1").
		Where("comment.project_id = ? AND comment.deleted_at IS NULL", projectID).
		Order("comment.created_at DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	summaries := make([]dto.CommentSummary, len(rows))
	for i, r := range rows {
		summaries[i] = dto.CommentSummary{
			ID:        r.ID,
			Nickname:  r.Nickname,
			Content:   r.Content,
			PostTitle: r.PostTitle,
			CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return summaries, nil
}

// PopularPosts 热门文章 top N（按浏览量）。
func (r *Repository) PopularPosts(ctx context.Context, projectID string, limit int) ([]dto.PostSummary, error) {
	projectID = normalizeProjectID(projectID)
	db := r.tx.DB(ctx)

	var rows []dto.PostSummary
	if err := db.Table("post").
		Select("id, title, status, created_at, view").
		Where("project_id = ? AND status = 2 AND deleted_at IS NULL", projectID).
		Order("view DESC").
		Limit(limit).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// PendingCounts 待处理计数。
func (r *Repository) PendingCounts(ctx context.Context, projectID string) (dto.Pending, error) {
	projectID = normalizeProjectID(projectID)
	db := r.tx.DB(ctx)

	var posts int64
	if err := db.Table("post").
		Where("project_id = ? AND status = 1 AND deleted_at IS NULL", projectID).
		Count(&posts).Error; err != nil {
		return dto.Pending{}, err
	}

	var comments int64
	if err := db.Table("comment").
		Where("project_id = ? AND status = 1 AND deleted_at IS NULL", projectID).
		Count(&comments).Error; err != nil {
		return dto.Pending{}, err
	}

	return dto.Pending{ReviewPosts: posts, ReviewComments: comments}, nil
}

// SeoAlerts 返回 SEO 提醒列表。
func (r *Repository) SeoAlerts(ctx context.Context, projectID string) []dto.SeoAlert {
	projectID = normalizeProjectID(projectID)
	db := r.tx.DB(ctx)
	var alerts []dto.SeoAlert

	var missingDesc int64
	if err := db.Table("post").
		Where("project_id = ? AND (description IS NULL OR description = '') AND deleted_at IS NULL", projectID).
		Count(&missingDesc).Error; err == nil && missingDesc > 0 {
		alerts = append(alerts, dto.SeoAlert{
			Type:    "warning",
			Message: fmt.Sprintf("%d 篇文章缺少 meta description", missingDesc),
		})
	}

	var missingAlt int64
	if err := db.Table("resource").
		Where("file_type = 1 AND (description IS NULL OR description = '') AND deleted_at IS NULL", projectID).
		Count(&missingAlt).Error; err == nil && missingAlt > 0 {
		alerts = append(alerts, dto.SeoAlert{
			Type:    "info",
			Message: fmt.Sprintf("%d 张图片未添加描述", missingAlt),
		})
	}

	if len(alerts) == 0 {
		alerts = append(alerts, dto.SeoAlert{
			Type:    "success",
			Message: "SEO 状态良好，未发现明显问题",
		})
	}
	return alerts
}

func normalizeProjectID(projectID string) string {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return defaultProjectID
	}
	return projectID
}
