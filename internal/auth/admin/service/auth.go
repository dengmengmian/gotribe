package service

// 本文件实现 Admin 认证业务逻辑，依赖 internal/auth/core 提供的 audience-aware Manager。

import (
	"context"
	"errors"
	"fmt"
	"time"

	adminrepo "gotribe/internal/admin/admin_user/repository"
	"gotribe/internal/auth/core"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	"gotribe/internal/model"
)

// Service 认证业务逻辑接口
type Service interface {
	Login(ctx context.Context, username, password, clientIP string) (*LoginResult, error)
	Refresh(ctx context.Context, userID int64, username string, issuedAt time.Time) (*LoginResult, error)
	Logout(ctx context.Context, userID int64) error
}

// LoginStage 表示登录流程当前所处阶段。
type LoginStage string

const (
	// LoginStageOK 表示登录已完成，Token 字段为可用的 access_token。
	LoginStageOK LoginStage = "ok"
	// LoginStageTOTPRequired 表示需要进行 TOTP 二次校验，StepToken 字段为短期凭证。
	LoginStageTOTPRequired LoginStage = "totp_required"
	// LoginStageBindRequired 表示策略要求强制绑定 TOTP（admin.totp.required=true 且账户未绑定），
	// 调用方需用 StepToken 走 /totp/enroll → /totp/enroll/confirm 完成首次绑定与登录。
	LoginStageBindRequired LoginStage = "bind_required"
)

// LoginResult 登录/刷新响应数据。字段按 Stage 取值不同而填充：
//   - Stage=ok          : Token / Expires / User 有效；可能附带 MFAReminder=true
//   - Stage=totp_required: StepToken / StepExpires 有效
type LoginResult struct {
	Stage       LoginStage
	Token       string
	Expires     time.Time
	StepToken   string
	StepExpires time.Time
	MFAReminder bool
	User        *model.Admin
}

// service 认证业务逻辑实现
type service struct {
	audience     string
	adminRepo    *adminrepo.Repository
	manager      *core.Manager
	tokenStore   *core.TokenStore
	lockout      *LockoutTracker
	totpService  TOTPService
	stepTokenTTL time.Duration
	totpRequired bool
}

// NewService 创建认证服务实例。audience 通常传 core.AudienceAdmin。
// tokenStore 用于登出吊销与刷新校验；可为 nil（Redis 不可用时降级为不吊销）。
// lockout 与 totpService 可为 nil（用于过渡期或单测），但建议生产环境注入。
// totpRequired=true 时，未绑 TOTP 的 admin 登录将返回 stage=bind_required 强制绑定。
func NewService(
	audience string,
	tx *database.TransactionManager,
	manager *core.Manager,
	tokenStore *core.TokenStore,
	lockout *LockoutTracker,
	totpService TOTPService,
	stepTokenTTL time.Duration,
	totpRequired bool,
) Service {
	return &service{
		audience:     audience,
		adminRepo:    adminrepo.NewRepository(tx),
		manager:      manager,
		tokenStore:   tokenStore,
		lockout:      lockout,
		totpService:  totpService,
		stepTokenTTL: stepTokenTTL,
		totpRequired: totpRequired,
	}
}

// Login 用户登录，验证身份并按状态决定是否签发 JWT 或要求 TOTP 二次校验。
// 流程：
//  1. 锁定预检（账户/IP 任一被锁则拒绝）
//  2. 密码校验失败 → 计数 + 返回错误
//  3. 密码校验成功 → 重置账户计数
//     a. 未启用 TOTP → 直接签发 access_token，附带 MFAReminder=true
//     b. 已启用 TOTP → 返回 step_token，stage=totp_required
func (s *service) Login(ctx context.Context, username, password, clientIP string) (*LoginResult, error) {
	if s.lockout != nil {
		if err := s.lockout.CheckBeforeLogin(ctx, username, clientIP); err != nil {
			var locked *LockedError
			if errors.As(err, &locked) {
				return nil, errs.AccountLocked("账户被临时锁定，请稍后再试", locked.LockedUntil.Unix(), locked.RemainingSeconds)
			}
			return nil, err
		}
	}

	u := &model.Admin{Username: username, Password: password}
	user, err := s.adminRepo.Login(ctx, u)
	if err != nil {
		if s.lockout != nil {
			_, locked := s.lockout.RecordFailure(ctx, username, clientIP)
			if locked != nil {
				return nil, errs.AccountLocked("失败次数过多，账户已被临时锁定", locked.LockedUntil.Unix(), locked.RemainingSeconds)
			}
		}
		return nil, err
	}

	if s.lockout != nil {
		s.lockout.Reset(ctx, username)
	}

	if s.totpService != nil {
		bound, err := s.totpService.IsBound(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		if bound {
			token, _, expires, err := s.manager.SignStepToken(user.ID, user.Username, core.StepTokenPurposeTOTPVerify, s.stepTokenTTL)
			if err != nil {
				return nil, errs.Internal("生成 step token 失败", err)
			}
			return &LoginResult{
				Stage:       LoginStageTOTPRequired,
				StepToken:   token,
				StepExpires: expires,
			}, nil
		}
		// 未绑定 + 策略强制 → 走「登录中途首次绑定」流程
		if s.totpRequired {
			token, _, expires, err := s.manager.SignStepToken(user.ID, user.Username, core.StepTokenPurposeTOTPBind, s.stepTokenTTL)
			if err != nil {
				return nil, errs.Internal("生成 step token 失败", err)
			}
			return &LoginResult{
				Stage:       LoginStageBindRequired,
				StepToken:   token,
				StepExpires: expires,
			}, nil
		}
	}

	subject := core.Subject{UserID: user.ID, Username: user.Username}
	token, expires, err := s.manager.SignAccessToken(s.audience, subject)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}
	return &LoginResult{
		Stage:       LoginStageOK,
		Token:       token,
		Expires:     expires,
		MFAReminder: s.totpService != nil, // 仅当系统启用了 TOTP 才提示绑定
		User:        user,
	}, nil
}

// Refresh 为已登录用户重新签发 JWT。
// issuedAt 为原 access token 的签发时间：若该会话已被登出吊销（invalid_before 之前签发），
// 则拒绝刷新，避免已登出 / 泄露的 token 被无限续期。
// Admin 身份无 project 维度，会话键以空 projectID 归一（与鉴权中间件一致）。
func (s *service) Refresh(ctx context.Context, userID int64, username string, issuedAt time.Time) (*LoginResult, error) {
	if s.tokenStore != nil {
		valid, err := s.tokenStore.IsAccessTokenValid(ctx, s.audience, "", userID, issuedAt)
		if err != nil {
			return nil, errs.ServiceUnavailable("会话状态校验暂不可用，请稍后重试", err)
		}
		if !valid {
			return nil, errs.Unauthorized("登录状态已失效，请重新登录")
		}
	}

	subject := core.Subject{
		UserID:   userID,
		Username: username,
	}
	token, expires, err := s.manager.SignAccessToken(s.audience, subject)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %w", err)
	}

	return &LoginResult{
		Stage:   LoginStageOK,
		Token:   token,
		Expires: expires,
	}, nil
}

// Logout 吊销该管理员当前 audience 下的所有会话：
// 设置 access token 失效时间为当前时刻，使此前签发的 access token 立即失效
// （由 JWTMiddleware 的 checker 与 Refresh 共同拒绝）。
// tokenStore 为 nil（Redis 不可用）时为无害的空操作。
func (s *service) Logout(ctx context.Context, userID int64) error {
	if s.tokenStore == nil {
		return nil
	}
	accessTTL, ok := s.manager.AccessTTL(s.audience)
	if !ok {
		return errs.Internal("未知的 audience，无法登出", nil)
	}
	return s.tokenStore.InvalidateUserSessions(ctx, s.audience, "", userID, time.Now().UTC(), accessTTL)
}
