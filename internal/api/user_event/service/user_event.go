// Package service implements user behavior event recording business logic.
package service

// 本文件实现用户行为上报的业务逻辑。

import (
	"context"
	"strconv"
	"strings"

	"gotribe/internal/core/errs"
	usereventdto "gotribe/internal/api/user_event/dto"
	"gotribe/internal/model"
	usereventrepo "gotribe/internal/api/user_event/repository"
)

// Service 负责封装用户行为事件相关的业务逻辑。
type Service struct {
	repo *usereventrepo.Repository
}

// NewService 创建用户行为事件服务实例。
func NewService(repo *usereventrepo.Repository) *Service {
	return &Service{repo: repo}
}

// Create 记录一条用户行为事件。
func (s *Service) Create(ctx context.Context, projectID string, userID int64, req usereventdto.CreateRequest, ip, userAgent string) error {
	if req.EventType == 0 {
		return errs.BadRequest("event_type is required", nil)
	}
	projectIDUint, _ := strconv.ParseInt(projectID, 10, 64)
	event := &model.UserEvent{
		UserID:      userID,
		ProjectID:   projectIDUint,
		EventType:   uint8(req.EventType),
		EventDetail: req.EventDetail,
		Duration:    req.Duration,
		IP:          ip,
		UserAgent:   userAgent,
		Referer:     strings.TrimSpace(req.Referer),
		Platform:    strings.TrimSpace(req.Platform),
	}
	if err := s.repo.Create(ctx, event); err != nil {
		return errs.Internal("create user event", err)
	}
	return nil
}
