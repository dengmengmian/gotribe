// Package service implements post listing, detail retrieval, caching, and tag assembly logic.
package service

// 本文件实现帖子读取、缓存和标签组合的业务逻辑。

import (
	"context"
	"crypto/subtle"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"
	"unicode"

	postdto "gotribe/internal/api/post/dto"
	postrepo "gotribe/internal/api/post/repository"
	postview "gotribe/internal/api/post/view"
	tagrepo "gotribe/internal/api/tag/repository"
	"gotribe/internal/core/cache"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	"gotribe/internal/model"
)

var htmlTagPattern = regexp.MustCompile(`<[^>]+>`)

// Service 负责封装文章相关的业务逻辑。
type Service struct {
	repo     *postrepo.Repository
	tags     *tagrepo.Repository
	cache    *cache.Store
	cacheTTL int
}

// NewService 创建文章服务实例。cacheTTL 由 config.Load 校验保证为正值。
func NewService(repo *postrepo.Repository, tags *tagrepo.Repository, cache *cache.Store, cacheTTL int) *Service {
	return &Service{
		repo:     repo,
		tags:     tags,
		cache:    cache,
		cacheTTL: cacheTTL,
	}
}

// List 查询文章列表并附带标签和分页信息。
func (s *Service) List(ctx context.Context, projectID string, query postdto.ListQuery) ([]postdto.PostResponse, database.Pagination, error) {
	var tagIDs []int64
	var err error
	if tag := strings.TrimSpace(query.Tag); tag != "" {
		tagIDs, err = s.tags.FindIDsByKeyword(ctx, tag)
		if err != nil {
			return nil, database.Pagination{}, errs.Internal("find tags", err)
		}
		if len(tagIDs) == 0 {
			page, perPage := database.NormalizePagination(query.Page, query.PerPage)
			return []postdto.PostResponse{}, database.Pagination{Page: page, PerPage: perPage, Total: 0}, nil
		}
	}

	filter := postrepo.ListFilter{
		Page:        query.Page,
		PerPage:     query.PerPage,
		Keyword:     query.Keyword,
		Status:      query.Status,
		Type:        query.Type,
		DynamicType: query.DynamicType,
		TagIDs:      tagIDs,
	}

	posts, total, err := s.repo.List(ctx, projectID, filter)
	if err != nil {
		return nil, database.Pagination{}, errs.Internal("list posts", err)
	}

	tagsByPostID, err := s.loadTagsByPostID(ctx, posts)
	if err != nil {
		return nil, database.Pagination{}, errs.Internal("list post tags", err)
	}

	page, perPage := database.NormalizePagination(query.Page, query.PerPage)
	items := make([]postdto.PostResponse, 0, len(posts))
	for _, item := range posts {
		items = append(items, toPostResponse(item, tagsByPostID[item.ID], false))
	}
	return items, database.Pagination{Page: page, PerPage: perPage, Total: total}, nil
}

// Detail 查询文章详情，并在需要时校验访问密码。
func (s *Service) Detail(ctx context.Context, projectID, postID, password string) (*postdto.PostResponse, error) {
	cacheKey := s.cache.PostDetailKey(projectID, postID)
	var cached postdto.PostResponse
	if ok, err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && ok {
		return &cached, nil
	}

	post, err := s.repo.GetByPostID(ctx, projectID, postID)
	if err != nil {
		return nil, errs.NotFound("post not found", err)
	}

	if post.IsPasswd != 0 && subtle.ConstantTimeCompare([]byte(post.PassWord), []byte(password)) != 1 {
		return nil, errs.Forbidden("post password is invalid")
	}

	tagsByPostID, err := s.loadTagsByPostID(ctx, []model.Post{*post})
	if err != nil {
		return nil, errs.Internal("load post tags", err)
	}

	resp := toPostResponse(*post, tagsByPostID[post.ID], true)
	if post.IsPasswd == 0 {
		_ = s.cache.SetJSON(ctx, cacheKey, resp, time.Duration(s.cacheTTL)*time.Minute)
	}
	return &resp, nil
}

// GetSummaries 按业务文章 ID 集合返回供其他模块复用的文章摘要。
func (s *Service) GetSummaries(ctx context.Context, projectID string, postIDs []string) (map[string]postview.Summary, error) {
	items, err := s.repo.ListByPostIDs(ctx, projectID, postIDs)
	if err != nil {
		return nil, errs.Internal("list posts by ids", err)
	}

	result := make(map[string]postview.Summary, len(items))
	for _, item := range items {
		result[item.PostID] = postview.Summary{
			PostID: item.PostID,
			Title:  item.Title,
			Type:   int16(item.Type),
			Status: int16(item.Status),
		}
	}
	return result, nil
}

// loadTagsByPostID 按文章集合批量加载标签关系。
func (s *Service) loadTagsByPostID(ctx context.Context, posts []model.Post) (map[int64][]model.Tag, error) {
	postIDs := make([]int64, 0, len(posts))
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}
	return s.tags.ListByPostIDs(ctx, postIDs)
}

// toPostResponse 将文章模型和标签信息组装为接口响应。
func toPostResponse(post model.Post, tagModels []model.Tag, includeBody bool) postdto.PostResponse {
	tags := make([]postdto.TagResponse, 0, len(tagModels))
	for _, item := range tagModels {
		tags = append(tags, postdto.TagResponse{
			ID:    int64(item.ID),
			Title: item.Title,
			Slug:  item.Slug,
			Color: item.Color,
		})
	}

	content := ""
	htmlContent := ""
	if includeBody {
		content = post.Content
		htmlContent = post.HtmlContent
	}
	return postdto.PostResponse{
		ID:           int64(post.ID),
		PostID:       post.PostID,
		Slug:         post.Slug,
		CategoryID:   int64(post.CategoryID),
		ProjectID:    fmt.Sprintf("%d", post.ProjectID),
		UserID:       int64(post.UserID),
		Author:       post.Author,
		Title:        post.Title,
		Content:      content,
		HTMLContent:  htmlContent,
		WordCount:    countHTMLWords(post.HtmlContent),
		Description:  post.Description,
		Icon:         post.Icon,
		View:         int64(post.View),
		Type:         int16(post.Type),
		Status:       int16(post.Status),
		UnitPrice:    int(post.UnitPrice),
		Location:     post.Location,
		People:       post.People,
		Time:         toTimeString(post.Time),
		ShowTime:     toTimeString(post.ShowTime),
		DynamicType:  post.DynamicType,
		Sort:         post.Sort,
		EventStartAt: toTimeString(post.EventStartAt),
		EventEndAt:   toTimeString(post.EventEndAt),
		RegisterURL:  post.RegisterURL,
		Tags:         tags,
		CreatedAt:    post.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    post.UpdatedAt.Format(time.RFC3339),
	}
}

// toTimeString 将时间指针转换为统一输出格式的字符串。
func toTimeString(value *time.Time) string {
	if value == nil || value.IsZero() {
		return ""
	}
	return value.Format(time.RFC3339)
}

// countHTMLWords 基于完整 HTML 正文估算可阅读文本字数。
func countHTMLWords(value string) int {
	text := html.UnescapeString(htmlTagPattern.ReplaceAllString(value, " "))
	count := 0
	inLatinWord := false
	for _, r := range text {
		if unicode.IsSpace(r) {
			inLatinWord = false
			continue
		}
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if r <= unicode.MaxASCII {
				if !inLatinWord {
					count++
					inLatinWord = true
				}
				continue
			}
			count++
			inLatinWord = false
			continue
		}
		inLatinWord = false
	}
	return count
}
