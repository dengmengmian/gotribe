// Package service implements login, token refresh, and logout business logic.
package service

// 本文件实现登录、刷新 token 和登出的业务逻辑。

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"gotribe/internal/auth/core"
	"gotribe/internal/auth/user/dto"
	"gotribe/internal/core/errs"
	authmodel "gotribe/internal/model"
)

type userReader interface {
	FindByIdentity(ctx context.Context, projectID, identity string) (*authmodel.AuthUser, error)
	FindByID(ctx context.Context, projectID string, userID int64) (*authmodel.AuthUser, error)
}

type refreshTokenStore interface {
	Save(ctx context.Context, audience, token string, session core.RefreshSession, ttl time.Duration) error
	Get(ctx context.Context, audience, token string) (core.RefreshSession, bool, error)
	Delete(ctx context.Context, audience, token string) error
	Rotate(ctx context.Context, audience, oldToken string, oldSession core.RefreshSession, newToken string, newSession core.RefreshSession, ttl time.Duration) (bool, error)
}

type tokenManager interface {
	SignAccessToken(audience string, subject core.Subject) (string, time.Time, error)
	GenerateRefreshToken() (string, error)
	AccessTTL(audience string) (time.Duration, bool)
	RefreshTTL(audience string) (time.Duration, bool)
}

// Service 负责封装认证相关的业务逻辑。
// 同一份实现服务于不同 audience（user / admin），通过构造时绑定 audience。
type Service struct {
	audience string
	users    userReader
	tokens   refreshTokenStore
	manager  tokenManager
}

// NewService 创建认证服务实例。audience 通过 core.AudienceUser / AudienceAdmin 等常量传入。
func NewService(audience string, users userReader, tokens refreshTokenStore, manager tokenManager) *Service {
	return &Service{
		audience: audience,
		users:    users,
		tokens:   tokens,
		manager:  manager,
	}
}

// dummyHash 用于时序攻击防护，确保无论用户是否存在，密码验证耗时一致。
const dummyHash = "$2a$10$abcdefghijklmnopqrstuuxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

// Login 校验身份信息并签发新的访问令牌和刷新令牌。
func (s *Service) Login(ctx context.Context, projectID string, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.users.FindByIdentity(ctx, projectID, strings.TrimSpace(req.Identity))
	if err != nil {
		// 用户不存在时仍执行一次 Verify 以消除时序攻击风险。
		_ = core.VerifyPassword(dummyHash, req.Password)
		return nil, errs.Unauthorized("invalid identity or password")
	}
	if user.Status == 0 {
		return nil, errs.Forbidden("user is disabled")
	}
	if !core.VerifyPassword(user.Password, req.Password) {
		return nil, errs.Unauthorized("invalid identity or password")
	}
	return s.issueTokens(ctx, user)
}

// Refresh 校验刷新令牌并重新签发会话凭证。
func (s *Service) Refresh(ctx context.Context, req dto.RefreshRequest) (*dto.AuthResponse, error) {
	session, ok, err := s.tokens.Get(ctx, s.audience, req.RefreshToken)
	if err != nil {
		return nil, errs.Internal("read refresh token", err)
	}
	if !ok {
		return nil, errs.Unauthorized("invalid refresh token")
	}

	user, err := s.users.FindByID(ctx, session.ProjectID, session.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Unauthorized("refresh token user not found")
		}
		return nil, errs.Internal("load refresh token user", err)
	}
	if user.Status == 0 {
		return nil, errs.Forbidden("user is disabled")
	}

	subject := core.Subject{UserID: user.ID, Username: user.Username, ProjectID: user.ProjectID}
	accessToken, _, err := s.manager.SignAccessToken(s.audience, subject)
	if err != nil {
		return nil, errs.Internal("generate access token", err)
	}

	refreshToken, err := s.manager.GenerateRefreshToken()
	if err != nil {
		return nil, errs.Internal("generate refresh token", err)
	}

	refreshTTL, ok := s.manager.RefreshTTL(s.audience)
	if !ok {
		return nil, errs.Internal("unknown audience", fmt.Errorf("audience %q not configured", s.audience))
	}
	accessTTL, _ := s.manager.AccessTTL(s.audience)

	newSession := core.RefreshSession{
		Audience:  s.audience,
		UserID:    user.ID,
		Username:  user.Username,
		ProjectID: user.ProjectID,
	}
	rotated, err := s.tokens.Rotate(ctx, s.audience, req.RefreshToken, session, refreshToken, newSession, refreshTTL)
	if err != nil {
		return nil, errs.Internal("rotate refresh token", err)
	}
	if !rotated {
		return nil, errs.Unauthorized("invalid refresh token")
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTTL.Seconds()),
		User: dto.UserSummary{
			ID:        user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			AvatarURL: user.AvatarURL,
			ProjectID: user.ProjectID,
		},
	}, nil
}

// Logout 使指定刷新令牌失效。
func (s *Service) Logout(ctx context.Context, currentUserID int64, req dto.LogoutRequest) error {
	session, ok, err := s.tokens.Get(ctx, s.audience, req.RefreshToken)
	if err != nil {
		return errs.Internal("read refresh token", err)
	}
	if !ok {
		return nil
	}
	if session.UserID != currentUserID {
		return errs.Forbidden("refresh token does not belong to current user")
	}
	if err := s.tokens.Delete(ctx, s.audience, req.RefreshToken); err != nil {
		return errs.Internal("delete refresh token", err)
	}
	return nil
}

// issueTokens 统一签发访问令牌和刷新令牌并组装响应。
func (s *Service) issueTokens(ctx context.Context, user *authmodel.AuthUser) (*dto.AuthResponse, error) {
	subject := core.Subject{UserID: user.ID, Username: user.Username, ProjectID: user.ProjectID}
	accessToken, _, err := s.manager.SignAccessToken(s.audience, subject)
	if err != nil {
		return nil, errs.Internal("generate access token", err)
	}

	refreshToken, err := s.manager.GenerateRefreshToken()
	if err != nil {
		return nil, errs.Internal("generate refresh token", err)
	}

	refreshTTL, ok := s.manager.RefreshTTL(s.audience)
	if !ok {
		return nil, errs.Internal("unknown audience", fmt.Errorf("audience %q not configured", s.audience))
	}
	accessTTL, _ := s.manager.AccessTTL(s.audience)

	if err := s.tokens.Save(ctx, s.audience, refreshToken, core.RefreshSession{
		Audience:  s.audience,
		UserID:    user.ID,
		Username:  user.Username,
		ProjectID: user.ProjectID,
	}, refreshTTL); err != nil {
		return nil, errs.Internal("store refresh token", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(accessTTL.Seconds()),
		User: dto.UserSummary{
			ID:        user.ID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Email:     user.Email,
			Phone:     user.Phone,
			AvatarURL: user.AvatarURL,
			ProjectID: user.ProjectID,
		},
	}, nil
}
