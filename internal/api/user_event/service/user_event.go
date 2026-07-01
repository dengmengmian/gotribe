// Package service implements user behavior event recording business logic.
package service

// 本文件实现用户行为上报的业务逻辑。

import (
	"context"
	"strconv"
	"strings"

	usereventdto "gotribe/internal/api/user_event/dto"
	usereventrepo "gotribe/internal/api/user_event/repository"
	"gotribe/internal/core/errs"
	"gotribe/internal/model"
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
	// 该接口在 public 组、project 来自 X-Project-ID 头，此前 ParseInt 静默吞错会把
	// 缺失/非法的 project 写成 0，污染数据。改为显式校验，拒绝空/非法/非正的 project。
	projectIDUint, err := strconv.ParseInt(strings.TrimSpace(projectID), 10, 64)
	if err != nil || projectIDUint <= 0 {
		return errs.BadRequest("invalid project", nil)
	}
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
