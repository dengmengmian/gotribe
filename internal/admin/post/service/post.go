package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.uber.org/zap"
	"gotribe/internal/admin/post/dto"
	postRepository "gotribe/internal/admin/post/repository"
	projectRepo "gotribe/internal/admin/project/repository"
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
}

// NewService 创建内容服务实例
func NewService(tx *database.TransactionManager, log *zap.SugaredLogger) Service {
	return &postService{
		postRepo:    postRepository.NewRepository(tx),
		projectRepo: projectRepo.NewRepository(tx),
		log:         log,
		tx:          tx,
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
	return s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.postRepo.Create(txCtx, &post); err != nil {
			if postRepository.IsInvalidTagIDError(err) {
				return errs.BadRequest(err.Error(), err)
			}
			return err
		}
		return nil
	})
}

// Update 更新内容
func (s *postService) Update(ctx context.Context, id int64, req *dto.UpdatePostRequest) error {
	oldPost, err := s.postRepo.Detail(ctx, id)
	if err != nil {
		return err
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
	return s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.postRepo.Update(txCtx, &oldPost); err != nil {
			if postRepository.IsInvalidTagIDError(err) {
				return errs.BadRequest(err.Error(), err)
			}
			return err
		}
		return nil
	})
}

// Delete 批量删除内容
func (s *postService) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		return s.postRepo.Delete(txCtx, ids)
	})
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
