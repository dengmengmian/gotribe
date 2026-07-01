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

	categoryrepo "gotribe/internal/api/category/repository"
	postdto "gotribe/internal/api/post/dto"
	postrepo "gotribe/internal/api/post/repository"
	postview "gotribe/internal/api/post/view"
	tagrepo "gotribe/internal/api/tag/repository"
	"gotribe/internal/core/cache"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	applog "gotribe/internal/core/logger"
	"gotribe/internal/model"

	"golang.org/x/sync/singleflight"
)

// listCacheEntry 用于文章列表缓存的序列化结构。
type listCacheEntry struct {
	Items []postdto.PostResponse `json:"items"`
	Meta  database.Pagination    `json:"meta"`
}

var htmlTagPattern = regexp.MustCompile(`<[^>]+>`)

// Service 负责封装文章相关的业务逻辑。
type Service struct {
	repo       *postrepo.Repository
	tags       *tagrepo.Repository
	categories *categoryrepo.Repository
	cache      *cache.Store
	cacheTTL   int
	// sf 合并同一 key 的并发回源，防止热点缓存过期瞬间大量请求同时打 DB（缓存击穿）。
	sf singleflight.Group
}

// NewService 创建文章服务实例。cacheTTL 由 config.Load 校验保证为正值。
func NewService(repo *postrepo.Repository, tags *tagrepo.Repository, categories *categoryrepo.Repository, cache *cache.Store, cacheTTL int) *Service {
	return &Service{
		repo:       repo,
		tags:       tags,
		categories: categories,
		cache:      cache,
		cacheTTL:   cacheTTL,
	}
}

// List 查询文章列表并附带标签和分页信息，结果按查询条件缓存 3 分钟。
func (s *Service) List(ctx context.Context, projectID string, query postdto.ListQuery) ([]postdto.PostResponse, database.Pagination, error) {
	cacheKey := s.cache.PostListKey(projectID, listCacheFilter(query))
	if entry, ok := s.readListCache(ctx, cacheKey); ok {
		return entry.Items, entry.Meta, nil
	}

	// singleflight 合并并发回源；返回值不缓存（如无匹配标签的空结果）时用 cached=false 标记。
	v, err, _ := s.sf.Do(cacheKey, func() (any, error) {
		if entry, ok := s.readListCache(ctx, cacheKey); ok {
			return entry, nil
		}

		var tagIDs []int64
		if tag := strings.TrimSpace(query.Tag); tag != "" {
			ids, err := s.tags.FindIDsByKeyword(ctx, tag)
			if err != nil {
				return nil, errs.Internal("find tags", err)
			}
			if len(ids) == 0 {
				page, perPage := database.NormalizePagination(query.Page, query.PerPage)
				return listCacheEntry{Items: []postdto.PostResponse{}, Meta: database.Pagination{Page: page, PerPage: perPage, Total: 0}}, nil
			}
			tagIDs = ids
		}

		filter := postrepo.ListFilter{
			Page:        query.Page,
			PerPage:     query.PerPage,
			Keyword:     query.Keyword,
			Status:      query.Status,
			Type:        query.Type,
			DynamicType: query.DynamicType,
			TagIDs:      tagIDs,
			CategoryID:  query.CategoryID,
		}

		posts, total, err := s.repo.List(ctx, projectID, filter)
		if err != nil {
			return nil, errs.Internal("list posts", err)
		}

		tagsByPostID, err := s.loadTagsByPostID(ctx, posts)
		if err != nil {
			return nil, errs.Internal("list post tags", err)
		}

		categoriesByID, err := s.loadCategoriesByID(ctx, posts)
		if err != nil {
			return nil, errs.Internal("list post categories", err)
		}

		page, perPage := database.NormalizePagination(query.Page, query.PerPage)
		items := make([]postdto.PostResponse, 0, len(posts))
		for _, item := range posts {
			items = append(items, toPostResponse(item, tagsByPostID[item.ID], categoriesByID[item.CategoryID], false))
		}

		entry := listCacheEntry{Items: items, Meta: database.Pagination{Page: page, PerPage: perPage, Total: total}}
		_ = s.cache.SetJSON(ctx, cacheKey, entry, cache.JitterTTL(time.Duration(s.cacheTTL)*time.Minute))
		return entry, nil
	})
	if err != nil {
		return nil, database.Pagination{}, err
	}
	entry := v.(listCacheEntry)
	return entry.Items, entry.Meta, nil
}

// readListCache 读取列表缓存；命中返回 (entry, true)。读错误（非 miss）记日志但按未命中处理。
func (s *Service) readListCache(ctx context.Context, cacheKey string) (listCacheEntry, bool) {
	var cached listCacheEntry
	ok, err := s.cache.GetJSON(ctx, cacheKey, &cached)
	if err != nil {
		applog.Warn(ctx, "post list cache read failed", "err", err)
		return listCacheEntry{}, false
	}
	return cached, ok
}

// detailResult 是 Detail 回源的中间结果：resp 为组装好的详情，
// isPasswd/password 用于回源后按「每个请求各自的密码」独立校验（密码文不进缓存、也不并发共享结果绕过校验）。
type detailResult struct {
	resp     *postdto.PostResponse
	isPasswd int64
	password string
}

// Detail 查询文章详情，并在需要时校验访问密码。
func (s *Service) Detail(ctx context.Context, projectID, postID, password string) (*postdto.PostResponse, error) {
	cacheKey := s.cache.PostDetailKey(projectID, postID)
	// 缓存里只存无密码文章，命中即可直接返回（无需再校验密码）。
	if resp, ok := s.readDetailCache(ctx, cacheKey); ok {
		return resp, nil
	}

	// singleflight 合并并发回源（DB 读取 + 组装），密码校验放在回源之后按请求各自进行。
	v, err, _ := s.sf.Do(cacheKey, func() (any, error) {
		if resp, ok := s.readDetailCache(ctx, cacheKey); ok {
			return &detailResult{resp: resp}, nil
		}

		post, err := s.repo.GetByPostID(ctx, projectID, postID)
		if err != nil {
			return nil, errs.NotFound("post not found", err)
		}

		tagsByPostID, err := s.loadTagsByPostID(ctx, []model.Post{*post})
		if err != nil {
			return nil, errs.Internal("load post tags", err)
		}
		categoriesByID, err := s.loadCategoriesByID(ctx, []model.Post{*post})
		if err != nil {
			return nil, errs.Internal("load post categories", err)
		}

		resp := toPostResponse(*post, tagsByPostID[post.ID], categoriesByID[post.CategoryID], true)
		if post.IsPasswd == 0 {
			_ = s.cache.SetJSON(ctx, cacheKey, resp, cache.JitterTTL(time.Duration(s.cacheTTL)*time.Minute))
		}
		return &detailResult{resp: &resp, isPasswd: post.IsPasswd, password: post.PassWord}, nil
	})
	if err != nil {
		return nil, err
	}

	result := v.(*detailResult)
	if result.isPasswd != 0 && subtle.ConstantTimeCompare([]byte(result.password), []byte(password)) != 1 {
		return nil, errs.Forbidden("post password is invalid")
	}
	return result.resp, nil
}

// readDetailCache 读取详情缓存；读错误（非 miss）记日志但按未命中处理。
func (s *Service) readDetailCache(ctx context.Context, cacheKey string) (*postdto.PostResponse, bool) {
	var cached postdto.PostResponse
	ok, err := s.cache.GetJSON(ctx, cacheKey, &cached)
	if err != nil {
		applog.Warn(ctx, "post detail cache read failed", "err", err)
		return nil, false
	}
	if !ok {
		return nil, false
	}
	return &cached, true
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

// listCacheFilter 把查询参数序列化为缓存 filter 字符串。
// 分页参数先经过 NormalizePagination 规范化，确保默认请求和显式传 page=1&per_page=10 共享同一缓存 key。
func listCacheFilter(query postdto.ListQuery) string {
	page, perPage := database.NormalizePagination(query.Page, query.PerPage)
	var sb strings.Builder
	fmt.Fprintf(&sb, "p%d|pp%d|k%s|t%s|dt%s", page, perPage, query.Keyword, query.Tag, query.DynamicType)
	if query.Status != nil {
		fmt.Fprintf(&sb, "|s%d", *query.Status)
	}
	if query.Type != nil {
		fmt.Fprintf(&sb, "|ty%d", *query.Type)
	}
	if query.CategoryID != nil {
		fmt.Fprintf(&sb, "|c%d", *query.CategoryID)
	}
	return sb.String()
}

// loadTagsByPostID 按文章集合批量加载标签关系。
func (s *Service) loadTagsByPostID(ctx context.Context, posts []model.Post) (map[int64][]model.Tag, error) {
	postIDs := make([]int64, 0, len(posts))
	for _, post := range posts {
		postIDs = append(postIDs, post.ID)
	}
	return s.tags.ListByPostIDs(ctx, postIDs)
}

// loadCategoriesByID 按文章集合批量加载分类信息。
func (s *Service) loadCategoriesByID(ctx context.Context, posts []model.Post) (map[int64]model.Category, error) {
	categoryIDs := make([]int64, 0, len(posts))
	seen := make(map[int64]struct{})
	for _, post := range posts {
		if post.CategoryID == 0 {
			continue
		}
		if _, ok := seen[post.CategoryID]; !ok {
			seen[post.CategoryID] = struct{}{}
			categoryIDs = append(categoryIDs, post.CategoryID)
		}
	}
	categories, err := s.categories.ListByIDs(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]model.Category, len(categories))
	for _, c := range categories {
		result[c.ID] = c
	}
	return result, nil
}

// toPostResponse 将文章模型和标签信息组装为接口响应。
func toPostResponse(post model.Post, tagModels []model.Tag, categoryModel model.Category, includeBody bool) postdto.PostResponse {
	tags := make([]postdto.TagResponse, 0, len(tagModels))
	for _, item := range tagModels {
		tags = append(tags, postdto.TagResponse{
			ID:    int64(item.ID),
			Title: item.Title,
			Slug:  item.Slug,
			Color: item.Color,
		})
	}

	var category *postdto.CategoryResponse
	if categoryModel.ID != 0 {
		category = &postdto.CategoryResponse{
			ID:    categoryModel.ID,
			Title: categoryModel.Title,
			Slug:  categoryModel.Slug,
		}
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
		Category:     category,
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
