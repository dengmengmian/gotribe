package service

import (
	"context"

	"gotribe/internal/admin/feedback/dto"
	feedbackRepository "gotribe/internal/admin/feedback/repository"
	"gotribe/internal/model"
	projectRepo "gotribe/internal/admin/project/repository"
	userRepo "gotribe/internal/admin/user/repository"

	"github.com/thoas/go-funk"
	"gotribe/internal/core/database"
)

// Service 反馈业务逻辑接口
type Service interface {
	List(ctx context.Context, req *dto.FeedbackListRequest) ([]*model.Feedback, int64, error)
	Delete(ctx context.Context, ids []int64) error
}

// service 反馈业务逻辑实现
type service struct {
	feedbackRepo *feedbackRepository.Repository
	userRepo     *userRepo.Repository
	projectRepo  *projectRepo.Repository
}

// NewFeedbackService 创建反馈服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		feedbackRepo: feedbackRepository.NewRepository(tx),
		userRepo:     userRepo.NewRepository(tx),
		projectRepo:  projectRepo.NewRepository(tx),
	}
}

// List 获取反馈列表
func (s *service) List(ctx context.Context, req *dto.FeedbackListRequest) ([]*model.Feedback, int64, error) {
	feedbacks, total, err := s.feedbackRepo.List(ctx, req)
	if err != nil {
		return feedbacks, total, err
	}
	feedbacks, err = s.getFeedbackOther(ctx, feedbacks)
	return feedbacks, total, err
}

func (s *service) getFeedbackOther(ctx context.Context, feedbacks []*model.Feedback) ([]*model.Feedback, error) {
	// 遍历feedback获取所有用户ID,查询用户信息并加进去
	userIdSet := make(map[int64]struct{})
	for _, feedback := range feedbacks {
		userIdSet[feedback.UserID] = struct{}{}
	}

	// 创建用户映射以提高查找效率
	userMap := make(map[int64]*model.User)
	for userID := range userIdSet {
		user, err := s.userRepo.Detail(ctx, userID)
		if err != nil {
			continue
		}
		userMap[userID] = &user
	}

	// 将用户信息附加到反馈中
	for _, feedback := range feedbacks {
		if user, ok := userMap[feedback.UserID]; ok {
			feedback.User = user
		}
	}

	// 追加项目信息
	projectIds := funk.UniqInt64(funk.Map(feedbacks, func(feedback *model.Feedback) int64 {
		return feedback.ProjectID
	}).([]int64))

	// 创建项目映射以提高查找效率
	projectMap := make(map[int64]*model.Project)
	for _, projectID := range projectIds {
		project, err := s.projectRepo.GetProjectByID(ctx, projectID)
		if err != nil {
			continue
		}
		projectMap[projectID] = &project
	}

	// 将项目信息附加到反馈中
	for _, feedback := range feedbacks {
		if project, ok := projectMap[feedback.ProjectID]; ok {
			feedback.Project = project
		}
	}

	return feedbacks, nil
}

// Delete 批量删除反馈。
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.feedbackRepo.Delete(ctx, ids)
}
