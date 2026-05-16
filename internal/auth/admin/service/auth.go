package service

// 本文件实现 Admin 认证业务逻辑，依赖 internal/auth/core 提供的 audience-aware Manager。

import (
	"context"
	"fmt"
	"time"

	adminrepo "gotribe/internal/admin/admin_user/repository"
	"gotribe/internal/auth/core"
	"gotribe/internal/core/database"
	"gotribe/internal/model"
)

// Service 认证业务逻辑接口
type Service interface {
	Login(ctx context.Context, username, password string) (*LoginResult, error)
	Refresh(ctx context.Context, userID int64, username string) (*LoginResult, error)
}

// LoginResult 登录/刷新响应数据
type LoginResult struct {
	Token   string
	Expires time.Time
	User    *model.Admin
}

// service 认证业务逻辑实现
type service struct {
	audience  string
	adminRepo *adminrepo.Repository
	manager   *core.Manager
}

// NewService 创建认证服务实例。audience 通常传 core.AudienceAdmin。
func NewService(audience string, tx *database.TransactionManager, manager *core.Manager) Service {
	return &service{
		audience:  audience,
		adminRepo: adminrepo.NewRepository(tx),
		manager:   manager,
	}
}

// Login 用户登录，验证身份并签发 JWT
func (s *service) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	u := &model.Admin{
		Username: username,
		Password: password,
	}

	user, err := s.adminRepo.Login(ctx, u)
	if err != nil {
		return nil, err
	}

	subject := core.Subject{
		UserID:   user.ID,
		Username: user.Username,
	}
	token, expires, err := s.manager.SignAccessToken(s.audience, subject)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	return &LoginResult{
		Token:   token,
		Expires: expires,
		User:    user,
	}, nil
}

// Refresh 为已登录用户重新签发 JWT
func (s *service) Refresh(ctx context.Context, userID int64, username string) (*LoginResult, error) {
	subject := core.Subject{
		UserID:   userID,
		Username: username,
	}
	token, expires, err := s.manager.SignAccessToken(s.audience, subject)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	return &LoginResult{
		Token:   token,
		Expires: expires,
	}, nil
}
