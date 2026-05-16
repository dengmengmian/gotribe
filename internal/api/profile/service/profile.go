// Package service implements user profile read, update, and password change logic.
package service

// 本文件实现当前用户资料读取、更新和改密码的业务逻辑。

import (
	"context"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"gotribe/internal/auth/core"
	"gotribe/internal/core/cache"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	applog "gotribe/internal/core/logger"
	profiledto "gotribe/internal/api/profile/dto"
	profilemodel "gotribe/internal/model"
	profilerepo "gotribe/internal/api/profile/repository"
	profileview "gotribe/internal/api/profile/view"
)

// sessionInvalidator 抽象按 audience 失效用户会话的能力。
type sessionInvalidator interface {
	InvalidateUserSessions(ctx context.Context, audience, projectID string, userID int64, invalidBefore time.Time, accessTTL time.Duration) error
}

// Service 负责封装当前用户资料相关的业务逻辑。
type Service struct {
	audience       string
	accessTokenTTL time.Duration
	repo           *profilerepo.Repository
	cache          *cache.Store
	tokens         sessionInvalidator
	tx             *database.TransactionManager
	cacheTTL       int
}

// NewService 创建当前用户资料服务实例。cacheTTL 由 config.Load 校验保证为正值。
// audience 与 accessTokenTTL 用于改密码后撤销该 audience 下的会话。
func NewService(
	audience string,
	accessTokenTTL time.Duration,
	repo *profilerepo.Repository,
	cache *cache.Store,
	tokens sessionInvalidator,
	tx *database.TransactionManager,
	cacheTTL int,
) *Service {
	return &Service{
		audience:       audience,
		accessTokenTTL: accessTokenTTL,
		repo:           repo,
		cache:          cache,
		tokens:         tokens,
		tx:             tx,
		cacheTTL:       cacheTTL,
	}
}

// GetMe 获取当前登录用户的资料信息。
func (s *Service) GetMe(ctx context.Context, projectID string, userID int64) (*profileview.MeView, error) {
	cacheKey := s.cache.ProfileKey(projectID, userID)
	var cached profileview.MeView
	if ok, err := s.cache.GetJSON(ctx, cacheKey, &cached); err == nil && ok {
		return &cached, nil
	}

	user, err := s.repo.GetByID(ctx, projectID, userID)
	if err != nil {
		return nil, errs.NotFound("user not found", err)
	}
	view := toMeView(user)
	_ = s.cache.SetJSON(ctx, cacheKey, view, time.Duration(s.cacheTTL)*time.Minute)
	return &view, nil
}

// UpdateMe 更新当前登录用户允许修改的资料字段。
func (s *Service) UpdateMe(ctx context.Context, projectID string, userID int64, req profiledto.UpdateProfileRequest) (*profileview.MeView, error) {
	updates := map[string]any{}
	if req.Nickname != nil {
		nickname := strings.TrimSpace(*req.Nickname)
		if nickname == "" {
			return nil, errs.BadRequest("nickname cannot be empty", nil)
		}
		if len(nickname) > 100 {
			return nil, errs.BadRequest("nickname is too long", nil)
		}
		updates["nickname"] = nickname
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if email != "" {
			if _, err := mail.ParseAddress(email); err != nil {
				return nil, errs.BadRequest("email format is invalid", err)
			}
		}
		updates["email"] = email
	}
	if req.Phone != nil {
		phone := strings.TrimSpace(*req.Phone)
		if phone != "" && !phonePattern.MatchString(phone) {
			return nil, errs.BadRequest("phone format is invalid", nil)
		}
		updates["phone"] = phone
	}
	if req.Sex != nil {
		sex := strings.TrimSpace(*req.Sex)
		if sex != "" && len(sex) != 1 {
			return nil, errs.BadRequest("sex must be a single character", nil)
		}
		updates["sex"] = sex
	}
	if req.Birthday != nil {
		birthday := strings.TrimSpace(*req.Birthday)
		if birthday == "" {
			updates["birthday"] = nil
		} else {
			parsed, err := time.Parse("2006-01-02", birthday)
			if err != nil {
				return nil, errs.BadRequest("birthday must use YYYY-MM-DD", err)
			}
			updates["birthday"] = parsed
		}
	}
	if req.Background != nil {
		updates["background"] = strings.TrimSpace(*req.Background)
	}
	if req.Ext != nil {
		updates["ext"] = strings.TrimSpace(*req.Ext)
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = strings.TrimSpace(*req.AvatarURL)
	}
	if len(updates) == 0 {
		return nil, errs.BadRequest("no fields to update", nil)
	}
	if err := s.repo.Update(ctx, projectID, userID, updates); err != nil {
		if errs.IsUniqueViolation(err) {
			switch errs.ConstraintName(err) {
			case "idx_user_email", "idx_user_project_email":
				return nil, errs.Conflict("email already exists", err)
			case "idx_user_phone", "idx_user_project_phone":
				return nil, errs.Conflict("phone already exists", err)
			default:
				return nil, errs.Conflict("profile field already exists", err)
			}
		}
		return nil, errs.Internal("update profile", err)
	}
	_ = s.cache.Delete(ctx, s.cache.ProfileKey(projectID, userID))
	return s.GetMe(ctx, projectID, userID)
}

// ChangePassword 校验旧密码并更新当前用户密码。
func (s *Service) ChangePassword(ctx context.Context, projectID string, userID int64, req profiledto.ChangePasswordRequest) error {
	if req.CurrentPassword == "" || req.NewPassword == "" {
		return errs.BadRequest("current_password and new_password are required", nil)
	}
	if req.CurrentPassword == req.NewPassword {
		return errs.BadRequest("new password must be different from current password", nil)
	}

	hashedPassword, err := core.HashPassword(req.NewPassword)
	if err != nil {
		return errs.BadRequest(err.Error(), err)
	}

	if err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		currentHash, err := s.repo.GetPasswordByID(txCtx, projectID, userID)
		if err != nil {
			return errs.NotFound("user not found", err)
		}
		if !core.VerifyPassword(currentHash, req.CurrentPassword) {
			return errs.Forbidden("current password is invalid")
		}
		if err := s.repo.UpdatePassword(txCtx, projectID, userID, hashedPassword); err != nil {
			return errs.Internal("update password", err)
		}
		return nil
	}); err != nil {
		return err
	}

	// 只有在所有会话都被成功失效后，才向客户端返回密码修改成功。
	if err := s.tokens.InvalidateUserSessions(ctx, s.audience, projectID, userID, time.Now().UTC(), s.accessTokenTTL); err != nil {
		applog.Error(ctx, "password updated but session invalidation failed", "err", err)
		return errs.Internal("invalidate user sessions", err)
	}
	return nil
}

// toMeView 将用户资料模型转换为 profile 模块内部视图。
func toMeView(user *profilemodel.UserProfile) profileview.MeView {
	birthday := ""
	if user.Birthday != nil && !user.Birthday.IsZero() {
		birthday = user.Birthday.Format("2006-01-02")
	}
	return profileview.MeView{
		ID:         int64(user.ID),
		Username:   user.Username,
		ProjectID:  user.ProjectID,
		Nickname:   user.Nickname,
		Email:      user.Email,
		Phone:      user.Phone,
		Sex:        user.Sex,
		Status:     user.Status,
		Birthday:   birthday,
		Background: user.Background,
		Ext:        user.Ext,
		AvatarURL:  user.AvatarURL,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  user.UpdatedAt.Format(time.RFC3339),
	}
}

var phonePattern = regexp.MustCompile(`^[0-9+\-]{6,20}$`)
