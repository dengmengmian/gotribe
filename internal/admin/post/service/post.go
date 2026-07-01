package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"gotribe/internal/admin/post/dto"
	postRepository "gotribe/internal/admin/post/repository"
	projectRepo "gotribe/internal/admin/project/repository"
	"gotribe/internal/core/cache"
	"gotribe/internal/core/constant"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/util"
	"gotribe/internal/model"
)

// Service 内容业务逻辑接口
type Service interface {
	Detail(ctx context.Context, id int64) (model.Post, error)
	List(ctx context.Context, req *dto.PostListRequest) ([]*model.Post, int64, error)
	Create(ctx context.Context, req *dto.CreatePostRequest) error
	Update(ctx context.Context, id int64, req *dto.UpdatePostRequest) error
	Delete(ctx context.Context, ids []int64) error
	Publish(ctx context.Context, id int64) error
}

// postService 内容业务逻辑实现
type postService struct {
	postRepo    *postRepository.Repository
	projectRepo *projectRepo.Repository
	log         *zap.SugaredLogger
	tx          *database.TransactionManager
	cache       *cache.Store
}

// NewService 创建内容服务实例
func NewService(tx *database.TransactionManager, log *zap.SugaredLogger, store *cache.Store) Service {
	return &postService{
		postRepo:    postRepository.NewRepository(tx),
		projectRepo: projectRepo.NewRepository(tx),
		log:         log,
		tx:          tx,
		cache:       store,
	}
}

// Detail 根据ID获取内容
func (s *postService) Detail(ctx context.Context, id int64) (model.Post, error) {
	return s.postRepo.Detail(ctx, id)
}

// List 获取内容列表
func (s *postService) List(ctx context.Context, req *dto.PostListRequest) ([]*model.Post, int64, error) {
	return s.postRepo.List(ctx, req)
}

// clearPostListCacheByProject 按项目精确清除文章列表缓存（best effort）。
func (s *postService) clearPostListCacheByProject(ctx context.Context, projectID int64) {
	if s.cache == nil || projectID <= 0 {
		return
	}
	pattern := s.cache.PostListPatternByProject(fmt.Sprintf("%d", projectID))
	if err := s.cache.DeleteByPattern(ctx, pattern); err != nil {
		s.log.Warnf("清除项目 %d 文章列表缓存失败: %v", projectID, err)
	}
}

// clearPostDetailCache 失效 ToC 文章详情缓存（best effort）。
// ToC 详情既可用 post_id 命中也可用 slug 命中，因此对同一篇文章需删除多个 key。
// identifiers 传入 post_id 及涉及的新旧 slug，空值会被忽略。
func (s *postService) clearPostDetailCache(ctx context.Context, projectID int64, identifiers ...string) {
	if s.cache == nil || projectID <= 0 {
		return
	}
	pid := fmt.Sprintf("%d", projectID)
	keys := make([]string, 0, len(identifiers))
	seen := make(map[string]struct{}, len(identifiers))
	for _, id := range identifiers {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		keys = append(keys, s.cache.PostDetailKey(pid, id))
	}
	if len(keys) == 0 {
		return
	}
	if err := s.cache.Delete(ctx, keys...); err != nil {
		s.log.Warnf("清除项目 %d 文章详情缓存失败: %v", projectID, err)
	}
}

// Create 创建内容
func (s *postService) Create(ctx context.Context, req *dto.CreatePostRequest) error {
	slug := req.Slug
	if slug == "" {
		slug = utils.GenerateSlug(req.Title)
	}
	imageStr := strings.Join(req.Images, ",")
	postTime, err := parseOptionalPostTime(req.Time, constant.TIME_FORMAT_SHORT)
	if err != nil {
		return err
	}
	showTime, err := parseOptionalPostTime(req.ShowTime, constant.TIME_FORMAT)
	if err != nil {
		return err
	}

	post := model.Post{
		Slug:        slug,
		CategoryID:  req.CategoryID,
		ProjectID:   req.ProjectID,
		UserID:      req.UserID,
		Author:      req.Author,
		Title:       req.Title,
		Content:     req.Content,
		HtmlContent: req.HtmlContent,
		Description: req.Description,
		Ext:         req.Ext,
		Tag:         req.Tag,
		Icon:        req.Icon,
		Type:        req.Type,
		IsTop:       req.IsTop,
		IsPasswd:    req.IsPasswd,
		ColumnID:    req.ColumnID,
		PassWord:    req.Password,
		Status:      req.Status,
		Time:        postTime,
		UnitPrice:   int64(utils.MoneyUtil.YuanToCents(req.UnitPrice)),
		People:      req.People,
		Location:    req.Location,
		Images:      imageStr,
		ShowTime:    showTime,
		Video:       req.Video,
	}
	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.postRepo.Create(txCtx, &post); err != nil {
			if postRepository.IsInvalidTagIDError(err) {
				return errs.BadRequest(err.Error(), err)
			}
			return err
		}
		return nil
	})
	if err == nil {
		s.clearPostListCacheByProject(ctx, post.ProjectID)
	}
	return err
}

// Update 更新内容
func (s *postService) Update(ctx context.Context, id int64, req *dto.UpdatePostRequest) error {
	oldPost, err := s.postRepo.Detail(ctx, id)
	if err != nil {
		return err
	}
	// 记录更新前的项目与 slug，便于失效可能因迁移项目 / 改 slug 而残留的旧缓存。
	originalProjectID := oldPost.ProjectID
	originalSlug := oldPost.Slug
	imageStr := strings.Join(req.Images, ",")
	postTime, err := parseOptionalPostTime(req.Time, constant.TIME_FORMAT_SHORT)
	if err != nil {
		return err
	}
	showTime, err := parseOptionalPostTime(req.ShowTime, constant.TIME_FORMAT)
	if err != nil {
		return err
	}
	slug := req.Slug
	if slug == "" {
		slug = utils.GenerateSlug(req.Title)
	}
	oldPost.Slug = slug
	oldPost.Title = req.Title
	oldPost.Description = req.Description
	oldPost.IsTop = req.IsTop
	oldPost.IsPasswd = req.IsPasswd
	oldPost.ProjectID = req.ProjectID
	oldPost.PassWord = req.Password
	oldPost.Type = req.Type
	oldPost.Icon = req.Icon
	oldPost.Ext = req.Ext
	oldPost.HtmlContent = req.HtmlContent
	oldPost.Content = req.Content
	oldPost.CategoryID = req.CategoryID
	oldPost.UserID = req.UserID
	oldPost.Author = req.Author
	oldPost.Status = req.Status
	oldPost.Tag = req.Tag
	oldPost.ColumnID = req.ColumnID
	oldPost.Time = postTime
	oldPost.UnitPrice = int64(utils.MoneyUtil.YuanToCents(req.UnitPrice))
	oldPost.People = req.People
	oldPost.Location = req.Location
	oldPost.Images = imageStr
	oldPost.Video = req.Video
	oldPost.ShowTime = showTime
	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.postRepo.Update(txCtx, &oldPost); err != nil {
			if postRepository.IsInvalidTagIDError(err) {
				return errs.BadRequest(err.Error(), err)
			}
			return err
		}
		return nil
	})
	if err == nil {
		// 原项目：列表 + 详情（含新旧 slug 与 post_id）都失效。
		s.clearPostListCacheByProject(ctx, originalProjectID)
		s.clearPostDetailCache(ctx, originalProjectID, oldPost.PostID, originalSlug, oldPost.Slug)
		// 若文章被迁移到新项目，新项目缓存也要失效。
		if oldPost.ProjectID != originalProjectID {
			s.clearPostListCacheByProject(ctx, oldPost.ProjectID)
			s.clearPostDetailCache(ctx, oldPost.ProjectID, oldPost.PostID, oldPost.Slug)
		}
	}
	return err
}

// Delete 批量删除内容
func (s *postService) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	// 删除前取出失效缓存所需字段；查询失败则不继续删除，避免留下无法清理的脏缓存。
	refs, err := s.postRepo.ListCacheRefsByIDs(ctx, ids)
	if err != nil {
		return errs.Internal("查询待删除文章缓存信息失败", err)
	}
	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.postRepo.Delete(txCtx, ids)
	})
	if err == nil {
		clearedList := make(map[int64]struct{}, len(refs))
		for _, ref := range refs {
			if _, ok := clearedList[ref.ProjectID]; !ok {
				clearedList[ref.ProjectID] = struct{}{}
				s.clearPostListCacheByProject(ctx, ref.ProjectID)
			}
			s.clearPostDetailCache(ctx, ref.ProjectID, ref.PostID, ref.Slug)
		}
	}
	return err
}

// Publish 发布内容
func (s *postService) Publish(ctx context.Context, id int64) error {
	oldPost, err := s.postRepo.Detail(ctx, id)
	if err != nil {
		return err
	}
	oldPost.Status = constant.POST_STATUS_PUBLIC
	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.postRepo.Update(txCtx, &oldPost)
	})
	if err != nil {
		return err
	}
	s.clearPostListCacheByProject(ctx, oldPost.ProjectID)
	s.clearPostDetailCache(ctx, oldPost.ProjectID, oldPost.PostID, oldPost.Slug)
	projectInfo, err := s.projectRepo.GetProjectByID(ctx, oldPost.ProjectID)
	if err != nil {
		s.log.Errorf("获取项目信息失败: %v", err)
		return nil
	}
	if !utils.IsEmpty(projectInfo.PushToken) {
		postURLWithID := projectInfo.PostURL + oldPost.Slug
		go func() {
			if _, err := utils.SEOUtil.PushBaidu(projectInfo.Domain, projectInfo.PushToken, postURLWithID); err != nil {
				s.log.Errorf("推送百度失败: %v", err)
			}
		}()
	}
	return nil
}

func parseOptionalPostTime(value string, primaryLayout string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}

	layouts := []string{primaryLayout}
	if primaryLayout == constant.TIME_FORMAT_SHORT {
		layouts = append(layouts, constant.TIME_FORMAT)
	}
	for _, layout := range layouts {
		parsed, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return &parsed, nil
		}
	}
	return nil, errors.New("invalid time format")
}
