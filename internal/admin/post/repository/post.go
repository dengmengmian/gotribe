package repository

import (
	"context"
	"errors"
	"fmt"
	"gotribe/internal/admin/post/dto"
	"gotribe/internal/core/database"
	"gotribe/internal/core/util"
	"gotribe/internal/model"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var ErrInvalidTagID = errors.New("无效的标签ID")

type Repository struct {
	tx *database.TransactionManager
}

func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

func buildPostOrder(req *dto.PostListRequest) string {
	sortByMap := map[string]string{
		"title":       "title",
		"author":      "author",
		"description": "description",
		"projectID":   "project_id",
		"project_id":  "project_id",
		"status":      "status",
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

// Detail 获取单个内容
func (r *Repository) Detail(ctx context.Context, id int64) (model.Post, error) {
	var post model.Post
	err := r.tx.DB(ctx).Where("id = ?", id).First(&post).Error
	if err != nil {
		return post, err
	}

	posts, err := r.GetPostOther(ctx, []*model.Post{&post})
	if err != nil {
		return post, err
	}
	if len(posts) == 0 {
		return post, nil
	}

	return *posts[0], nil
}

// List 获取内容列表
func (r *Repository) List(ctx context.Context, req *dto.PostListRequest) ([]*model.Post, int64, error) {
	var list []*model.Post
	db := r.tx.DB(ctx).Model(&model.Post{})

	title := strings.TrimSpace(req.Title)
	if !utils.IsEmpty(title) {
		db = db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", title))
	}
	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	}
	if req.ProjectID > 0 {
		db = db.Where("project_id = ?", req.ProjectID)
	}
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	db = db.Order(buildPostOrder(req))

	page, perPage := database.NormalizePagination(int(req.PageNum), int(req.PageSize))
	err = db.Offset((page - 1) * perPage).Limit(perPage).Find(&list).Error

	if err != nil {
		return nil, 0, err
	}
	// 调用 GetPostOther 并处理返回值
	list, err = r.GetPostOther(ctx, list)
	return list, total, err
}

func (r *Repository) GetPostOther(ctx context.Context, posts []*model.Post) ([]*model.Post, error) {
	// 收集所有需要查询的 CategoryID, Tag, ProjectID
	categoryIDSet := make(map[int64]struct{})
	projectIDSet := make(map[int64]struct{})
	postIDs := make([]int64, 0, len(posts))

	for _, m := range posts {
		if m.CategoryID > 0 {
			categoryIDSet[m.CategoryID] = struct{}{}
		}
		if m.ProjectID > 0 {
			projectIDSet[m.ProjectID] = struct{}{}
		}
		if m.ID > 0 {
			postIDs = append(postIDs, m.ID)
		}
	}

	// 转换为切片
	categoryIDs := make([]int64, 0, len(categoryIDSet))
	for id := range categoryIDSet {
		categoryIDs = append(categoryIDs, id)
	}

	projectIDs := make([]int64, 0, len(projectIDSet))
	for id := range projectIDSet {
		projectIDs = append(projectIDs, id)
	}

	// 批量查询 Category
	var categories []*model.Category
	if len(categoryIDs) > 0 {
		if err := r.tx.DB(ctx).Where("id IN (?)", categoryIDs).Find(&categories).Error; err != nil {
			return nil, err
		}
	}
	// 批量查询 PostTag 关联
	postTagMap := make(map[int64][]int64)
	tagIDSet := make(map[int64]struct{})
	if len(postIDs) > 0 {
		var postTags []model.PostTag
		if err := r.tx.DB(ctx).Where("post_id IN (?)", postIDs).Find(&postTags).Error; err != nil {
			return nil, err
		}
		for _, postTag := range postTags {
			postTagMap[postTag.PostID] = append(postTagMap[postTag.PostID], postTag.TagID)
			tagIDSet[postTag.TagID] = struct{}{}
		}
	}

	// 批量查询 Tag
	var allTags []*model.Tag
	tagIDs := make([]int64, 0, len(tagIDSet))
	for tagID := range tagIDSet {
		tagIDs = append(tagIDs, tagID)
	}
	if len(tagIDs) > 0 {
		if err := r.tx.DB(ctx).Where("id IN (?)", tagIDs).Find(&allTags).Error; err != nil {
			return nil, err
		}
	}

	// 批量查询 Project
	var projects []*model.Project
	if len(projectIDs) > 0 {
		if err := r.tx.DB(ctx).Where("id IN (?)", projectIDs).Find(&projects).Error; err != nil {
			return nil, err
		}
	}

	// 将查询结果赋值给 posts
	categoryMap := make(map[int64]*model.Category)
	for _, category := range categories {
		categoryMap[category.ID] = category
	}

	tagMap := make(map[int64]*model.Tag)
	for _, tag := range allTags {
		tagMap[tag.ID] = tag
	}

	projectMap := make(map[int64]*model.Project)
	for _, project := range projects {
		projectMap[project.ID] = project
	}

	for _, m := range posts {
		if category, ok := categoryMap[m.CategoryID]; ok {
			m.Category = category
		}
		var tags []*model.Tag
		for _, tagID := range postTagMap[m.ID] {
			if tag, ok := tagMap[tagID]; ok {
				tags = append(tags, tag)
			}
		}
		m.Tags = tags
		m.Tag = formatPostTagIDs(tags)
		if project, ok := projectMap[m.ProjectID]; ok {
			m.Project = project
		}
	}
	return posts, nil
}

// Create 创建内容
func (r *Repository) Create(ctx context.Context, post *model.Post) error {
	db := r.tx.DB(ctx)
	if err := db.Create(post).Error; err != nil {
		return err
	}
	return syncPostTags(db, post.ID, post.Tag)
}

// Update 更新内容
func (r *Repository) Update(ctx context.Context, post *model.Post) error {
	db := r.tx.DB(ctx)
	if err := db.Model(&model.Post{}).Where("id = ?", post.ID).Updates(buildPostUpdateMap(post)).Error; err != nil {
		return err
	}
	return syncPostTags(db, post.ID, post.Tag)
}

// Delete 批量删除
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	var posts []model.Post
	for _, id := range ids {
		// 根据ID获取标签
		post, err := r.Detail(ctx, id)
		if err != nil {
			return fmt.Errorf("未获取到ID为%d的内容", id)
		}
		posts = append(posts, post)
	}

	db := r.tx.DB(ctx)
	postIDs := make([]int64, 0, len(posts))
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}
	if len(postIDs) > 0 {
		if err := db.Where("post_id IN (?)", postIDs).Delete(&model.PostTag{}).Error; err != nil {
			return err
		}
	}
	return db.Delete(&posts).Error
}

func syncPostTags(tx *gorm.DB, postID int64, rawTagIDs string) error {
	tagIDs, err := parsePostTagIDs(rawTagIDs)
	if err != nil {
		return err
	}

	if err := tx.Where("post_id = ?", postID).Delete(&model.PostTag{}).Error; err != nil {
		return err
	}
	if len(tagIDs) == 0 {
		return nil
	}

	postTags := make([]model.PostTag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		postTags = append(postTags, model.PostTag{
			PostID: postID,
			TagID:  tagID,
		})
	}

	return tx.Create(&postTags).Error
}

func parsePostTagIDs(raw string) ([]int64, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	values := strings.Split(raw, ",")
	seen := make(map[int64]struct{}, len(values))
	tagIDs := make([]int64, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		tagID, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidTagID, value)
		}

		if _, ok := seen[tagID]; ok {
			continue
		}
		seen[tagID] = struct{}{}
		tagIDs = append(tagIDs, tagID)
	}

	sort.Slice(tagIDs, func(i, j int) bool {
		return tagIDs[i] < tagIDs[j]
	})
	return tagIDs, nil
}

func IsInvalidTagIDError(err error) bool {
	return errors.Is(err, ErrInvalidTagID)
}

func formatPostTagIDs(tags []*model.Tag) string {
	if len(tags) == 0 {
		return ""
	}

	tagIDs := make([]string, 0, len(tags))
	for _, tag := range tags {
		if tag == nil {
			continue
		}
		tagIDs = append(tagIDs, strconv.FormatUint(uint64(tag.ID), 10))
	}
	return strings.Join(tagIDs, ",")
}

func buildPostUpdateMap(post *model.Post) map[string]interface{} {
	return map[string]interface{}{
		"category_id":  post.CategoryID,
		"project_id":   post.ProjectID,
		"column_id":    post.ColumnID,
		"user_id":      post.UserID,
		"author":       post.Author,
		"title":        post.Title,
		"slug":         post.Slug,
		"content":      post.Content,
		"html_content": post.HtmlContent,
		"description":  post.Description,
		"ext":          post.Ext,
		"icon":         post.Icon,
		"view":         post.View,
		"type":         post.Type,
		"is_top":       post.IsTop,
		"is_passwd":    post.IsPasswd,
		"pass_word":    post.PassWord,
		"status":       post.Status,
		"unit_price":   post.UnitPrice,
		"location":     post.Location,
		"people":       post.People,
		"time":         post.Time,
		"images":       post.Images,
		"show_time":    post.ShowTime,
		"video":        post.Video,
	}
}
