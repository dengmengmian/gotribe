package service

// 本文件实现 admin TOTP 二次校验业务逻辑：
//   - Bind：生成新 secret + otpauth URL + 备份码（不立即启用）
//   - Confirm：用一次 6 位码确认绑定并启用
//   - Verify：登录后两步校验，支持 TOTP code 或备份码
//   - Delete：自助解绑（需当前 6 位码）
//   - Reset：超管强制重置
// secret 用 AES-256-GCM 加密后存数据库；备份码 bcrypt hash 存数据库（明文仅在 Bind 时一次性返回）。

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	adminrepo "gotribe/internal/admin/admin_user/repository"
	authrepo "gotribe/internal/auth/admin/repository"
	"gotribe/internal/auth/core"
	"gotribe/internal/core/cache"
	"gotribe/internal/core/config"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	applog "gotribe/internal/core/logger"
	utils "gotribe/internal/core/util"
	"gotribe/internal/model"

	"github.com/redis/go-redis/v9"
)

// TOTPService 提供 admin TOTP 二次校验业务能力。
type TOTPService interface {
	Status(ctx context.Context, adminID int64) (*TOTPStatus, error)
	Bind(ctx context.Context, admin *model.Admin) (*TOTPBindResult, error)
	Confirm(ctx context.Context, adminID int64, code string) error
	VerifyAndIssue(ctx context.Context, stepToken, code string) (*LoginResult, error)
	Delete(ctx context.Context, admin *model.Admin, currentCode string) error
	AdminReset(ctx context.Context, targetAdminID int64) error
	IsBound(ctx context.Context, adminID int64) (bool, error)
	// EnrollPending：admin.totp.required=true 时，用 step_token(purpose=totp_bind) 发起首次绑定，
	// 返回 secret/QR/备份码。step_token 不会被消费，留给后续 ConfirmEnrollPending 使用。
	EnrollPending(ctx context.Context, stepToken string) (*TOTPBindResult, error)
	// ConfirmEnrollPending：用同一 step_token + 6 位码完成首次绑定并签发 access_token；
	// 成功后该 step_token 失效（jti 黑名单）。
	ConfirmEnrollPending(ctx context.Context, stepToken, code string) (*LoginResult, error)
}

// TOTPStatus 描述当前账户的 TOTP 绑定状态。
type TOTPStatus struct {
	Bound                  bool       `json:"bound"`
	Enabled                bool       `json:"enabled"`
	LastUsedAt             *time.Time `json:"last_used_at,omitempty"`
	RemainingRecoveryCodes int        `json:"remaining_recovery_codes"`
}

// TOTPBindResult 是 Bind 阶段一次性返回的内容（secret 明文 + URL + 备份码）。
type TOTPBindResult struct {
	Secret         string   `json:"secret"`           // base32 字符串，用户可手动输入
	OTPAuthURL     string   `json:"otpauth_url"`      // otpauth://totp/... 用于生成 QR
	RecoveryCodes  []string `json:"recovery_codes"`   // 10 个明文，仅本次返回
}

// recoveryCodeRecord 内部存储结构。
type recoveryCodeRecord struct {
	Hash   string `json:"hash"`
	UsedAt *int64 `json:"used_at,omitempty"`
}

// totpService 实现。
type totpService struct {
	cfg       config.AdminTOTPConfig
	repo      *authrepo.TOTPRepository
	adminRepo *adminrepo.Repository
	cipher    *utils.AESGCMCipher
	manager   *core.Manager
	redis     redis.UniversalClient
	keys      *cache.KeyBuilder
	audience  string
}

// NewTOTPService 构造 TOTPService。
func NewTOTPService(
	cfg config.AdminTOTPConfig,
	audience string,
	tx *database.TransactionManager,
	manager *core.Manager,
	redisClient redis.UniversalClient,
	keys *cache.KeyBuilder,
) (TOTPService, error) {
	cipher, err := utils.NewAESGCMCipher([]byte(cfg.SecretEncryptionKey))
	if err != nil {
		return nil, fmt.Errorf("init totp cipher: %w", err)
	}
	return &totpService{
		cfg:       cfg,
		repo:      authrepo.NewTOTPRepository(tx),
		adminRepo: adminrepo.NewRepository(tx),
		cipher:    cipher,
		manager:   manager,
		redis:     redisClient,
		keys:      keys,
		audience:  audience,
	}, nil
}

// IsBound 判断 admin 是否完成绑定（enabled=true）。
func (s *totpService) IsBound(ctx context.Context, adminID int64) (bool, error) {
	record, err := s.repo.GetByAdminID(ctx, adminID)
	if err != nil {
		return false, err
	}
	return record != nil && record.Enabled, nil
}

// Status 返回 TOTP 状态。
func (s *totpService) Status(ctx context.Context, adminID int64) (*TOTPStatus, error) {
	record, err := s.repo.GetByAdminID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	st := &TOTPStatus{}
	if record == nil {
		return st, nil
	}
	st.Bound = true
	st.Enabled = record.Enabled
	if record.LastUsedAt != nil {
		ts := time.Unix(*record.LastUsedAt, 0)
		st.LastUsedAt = &ts
	}
	codes, _ := s.decodeRecoveryCodes(record.RecoveryCodes)
	for _, code := range codes {
		if code.UsedAt == nil {
			st.RemainingRecoveryCodes++
		}
	}
	return st, nil
}

// Bind 生成新 secret + 备份码。允许重新绑定（覆盖未启用的记录）。
// 若已启用，需先解绑。
func (s *totpService) Bind(ctx context.Context, admin *model.Admin) (*TOTPBindResult, error) {
	existing, err := s.repo.GetByAdminID(ctx, admin.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Enabled {
		return nil, errs.TOTPAlreadyBound("已绑定 TOTP，如需重新绑定请先解绑")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.cfg.Issuer,
		AccountName: admin.Username,
		Period:      uint(s.cfg.Period),
		Digits:      otp.Digits(s.cfg.Digits),
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, errs.Internal("生成 TOTP secret 失败", err)
	}

	secretCipher, err := s.cipher.EncryptToString(key.Secret())
	if err != nil {
		return nil, errs.Internal("加密 TOTP secret 失败", err)
	}

	recoveryCodes, recoveryRecords, err := s.generateRecoveryCodes()
	if err != nil {
		return nil, err
	}
	recoveryJSON, err := json.Marshal(recoveryRecords)
	if err != nil {
		return nil, errs.Internal("序列化备份码失败", err)
	}

	record := &model.AdminTOTP{
		AdminID:       admin.ID,
		SecretCipher:  secretCipher,
		Enabled:       false,
		RecoveryCodes: string(recoveryJSON),
	}
	if err := s.repo.Upsert(ctx, record); err != nil {
		return nil, err
	}

	return &TOTPBindResult{
		Secret:        key.Secret(),
		OTPAuthURL:    key.URL(),
		RecoveryCodes: recoveryCodes,
	}, nil
}

// Confirm 用一次 6 位码确认绑定并启用记录。
func (s *totpService) Confirm(ctx context.Context, adminID int64, code string) error {
	record, err := s.repo.GetByAdminID(ctx, adminID)
	if err != nil {
		return err
	}
	if record == nil {
		return errs.TOTPNotBound("尚未发起绑定流程")
	}
	if record.Enabled {
		return errs.TOTPAlreadyBound("已启用，无需重复确认")
	}

	secret, err := s.cipher.DecryptFromString(record.SecretCipher)
	if err != nil {
		return errs.Internal("解密 TOTP secret 失败", err)
	}
	if !s.validateCode(secret, code) {
		return errs.TOTPInvalid("验证码错误")
	}

	if err := s.repo.MarkEnabled(ctx, adminID); err != nil {
		return err
	}
	now := time.Now().Unix()
	_ = s.repo.UpdateLastUsedAt(ctx, adminID, now)
	return nil
}

// VerifyAndIssue 校验 step_token + TOTP code（或备份码），通过后签发 access_token。
func (s *totpService) VerifyAndIssue(ctx context.Context, stepToken, code string) (*LoginResult, error) {
	claims, err := s.manager.VerifyStepToken(stepToken, core.StepTokenPurposeTOTPVerify)
	if err != nil {
		return nil, errs.Unauthorized("step token 无效或已过期")
	}

	if used, _ := s.isStepJTIUsed(ctx, claims.ID); used {
		return nil, errs.Unauthorized("step token 已被使用")
	}

	record, err := s.repo.GetByAdminID(ctx, claims.AdminID)
	if err != nil {
		return nil, err
	}
	if record == nil || !record.Enabled {
		return nil, errs.TOTPNotBound("账户未启用 TOTP")
	}

	secret, err := s.cipher.DecryptFromString(record.SecretCipher)
	if err != nil {
		return nil, errs.Internal("解密 TOTP secret 失败", err)
	}

	verified := false
	if s.validateCode(secret, code) {
		verified = true
	} else if matched, updated, err := s.consumeRecoveryCode(record.RecoveryCodes, code); err != nil {
		return nil, err
	} else if matched {
		verified = true
		if err := s.repo.UpdateRecoveryCodes(ctx, claims.AdminID, updated); err != nil {
			return nil, err
		}
	}
	if !verified {
		return nil, errs.TOTPInvalid("验证码错误")
	}

	now := time.Now().Unix()
	_ = s.repo.UpdateLastUsedAt(ctx, claims.AdminID, now)
	s.markStepJTIUsed(ctx, claims.ID, time.Until(claims.ExpiresAt.Time))

	admin, err := s.adminRepo.GetAdminByID(ctx, claims.AdminID)
	if err != nil {
		return nil, errs.Unauthorized("用户不存在")
	}

	subject := core.Subject{UserID: admin.ID, Username: admin.Username}
	token, expires, err := s.manager.SignAccessToken(s.audience, subject)
	if err != nil {
		return nil, errs.Internal("生成 token 失败", err)
	}
	return &LoginResult{
		Stage:   LoginStageOK,
		Token:   token,
		Expires: expires,
		User:    &admin,
	}, nil
}

// Delete 自助解绑（需当前 6 位码）。
func (s *totpService) Delete(ctx context.Context, admin *model.Admin, currentCode string) error {
	record, err := s.repo.GetByAdminID(ctx, admin.ID)
	if err != nil {
		return err
	}
	if record == nil || !record.Enabled {
		return errs.TOTPNotBound("尚未启用 TOTP")
	}
	secret, err := s.cipher.DecryptFromString(record.SecretCipher)
	if err != nil {
		return errs.Internal("解密 TOTP secret 失败", err)
	}
	if !s.validateCode(secret, currentCode) {
		return errs.TOTPInvalid("验证码错误")
	}
	return s.repo.DeleteByAdminID(ctx, admin.ID)
}

// AdminReset 超管强制重置他人 TOTP（不校验 code）。
func (s *totpService) AdminReset(ctx context.Context, targetAdminID int64) error {
	return s.repo.DeleteByAdminID(ctx, targetAdminID)
}

// validateCode 用 pquerna/otp 校验 code，支持时钟漂移 ±skew 步。
func (s *totpService) validateCode(secret, code string) bool {
	code = strings.TrimSpace(code)
	if len(code) != s.cfg.Digits {
		return false
	}
	valid, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    uint(s.cfg.Period),
		Skew:      uint(s.cfg.Skew),
		Digits:    otp.Digits(s.cfg.Digits),
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return false
	}
	return valid
}

// generateRecoveryCodes 生成 N 个明文备份码及其 hash 表，hash 用 bcrypt cost=10。
func (s *totpService) generateRecoveryCodes() ([]string, []recoveryCodeRecord, error) {
	count := s.cfg.RecoveryCodesCount
	plain := make([]string, 0, count)
	records := make([]recoveryCodeRecord, 0, count)
	for i := 0; i < count; i++ {
		code, err := randomRecoveryCode()
		if err != nil {
			return nil, nil, errs.Internal("生成备份码失败", err)
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			return nil, nil, errs.Internal("hash 备份码失败", err)
		}
		plain = append(plain, code)
		records = append(records, recoveryCodeRecord{Hash: string(hash)})
	}
	return plain, records, nil
}

// consumeRecoveryCode 尝试用 input 匹配未使用的备份码。匹配成功返回更新后的 JSON 字符串。
func (s *totpService) consumeRecoveryCode(stored, input string) (bool, string, error) {
	input = strings.ReplaceAll(strings.TrimSpace(input), "-", "")
	if input == "" {
		return false, "", nil
	}
	records, err := s.decodeRecoveryCodes(stored)
	if err != nil {
		return false, "", err
	}
	for i := range records {
		if records[i].UsedAt != nil {
			continue
		}
		if err := bcrypt.CompareHashAndPassword([]byte(records[i].Hash), []byte(input)); err == nil {
			now := time.Now().Unix()
			records[i].UsedAt = &now
			updated, mErr := json.Marshal(records)
			if mErr != nil {
				return true, "", errs.Internal("序列化备份码失败", mErr)
			}
			return true, string(updated), nil
		}
	}
	return false, "", nil
}

func (s *totpService) decodeRecoveryCodes(raw string) ([]recoveryCodeRecord, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var records []recoveryCodeRecord
	if err := json.Unmarshal([]byte(raw), &records); err != nil {
		return nil, errs.Internal("解析备份码失败", err)
	}
	return records, nil
}

// isStepJTIUsed 检查 step_token jti 是否已用过。Redis 故障 fail-open（视为未用）+ 日志。
func (s *totpService) isStepJTIUsed(ctx context.Context, jti string) (bool, error) {
	if s.redis == nil || s.keys == nil {
		return false, nil
	}
	key := s.keys.AdminTOTPStepJTIKey(jti)
	_, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		applog.Sugared().Warnf("TOTPService: 读取 step jti 失败 jti=%s: %v", jti, err)
		return false, nil
	}
	return true, nil
}

// markStepJTIUsed 将 jti 写入黑名单，TTL 与 token 剩余有效期一致。
func (s *totpService) markStepJTIUsed(ctx context.Context, jti string, ttl time.Duration) {
	if s.redis == nil || s.keys == nil || ttl <= 0 {
		return
	}
	key := s.keys.AdminTOTPStepJTIKey(jti)
	if err := s.redis.Set(ctx, key, "1", ttl).Err(); err != nil {
		applog.Sugared().Warnf("TOTPService: 写入 step jti 失败 jti=%s: %v", jti, err)
	}
}

// randomRecoveryCode 生成形如 "XXXX-XXXX" 的 8 字符备份码（base32 字符集去除易混字符）。
// 实际产出存储时去掉连字符，比较时也去掉。
func randomRecoveryCode() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 去除 0/O/1/I
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i := range buf {
		buf[i] = alphabet[int(buf[i])%len(alphabet)]
	}
	return string(buf[:4]) + "-" + string(buf[4:]), nil
}

// EnrollPending：用 step_token(purpose=totp_bind) 发起首次绑定。step_token 暂不消费，
// 让用户在同一 step 内可以反复看 secret 和最终提交确认；ConfirmEnrollPending 才会作废 jti。
func (s *totpService) EnrollPending(ctx context.Context, stepToken string) (*TOTPBindResult, error) {
	claims, err := s.manager.VerifyStepToken(stepToken, core.StepTokenPurposeTOTPBind)
	if err != nil {
		return nil, errs.Unauthorized("step token 无效或已过期")
	}
	if used, _ := s.isStepJTIUsed(ctx, claims.ID); used {
		return nil, errs.Unauthorized("step token 已被使用")
	}

	existing, err := s.repo.GetByAdminID(ctx, claims.AdminID)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Enabled {
		return nil, errs.TOTPAlreadyBound("已绑定 TOTP，无需再次绑定")
	}

	admin, err := s.adminRepo.GetAdminByID(ctx, claims.AdminID)
	if err != nil {
		return nil, errs.Unauthorized("用户不存在")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.cfg.Issuer,
		AccountName: admin.Username,
		Period:      uint(s.cfg.Period),
		Digits:      otp.Digits(s.cfg.Digits),
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, errs.Internal("生成 TOTP secret 失败", err)
	}
	secretCipher, err := s.cipher.EncryptToString(key.Secret())
	if err != nil {
		return nil, errs.Internal("加密 TOTP secret 失败", err)
	}
	recoveryCodes, recoveryRecords, err := s.generateRecoveryCodes()
	if err != nil {
		return nil, err
	}
	recoveryJSON, err := json.Marshal(recoveryRecords)
	if err != nil {
		return nil, errs.Internal("序列化备份码失败", err)
	}
	record := &model.AdminTOTP{
		AdminID:       claims.AdminID,
		SecretCipher:  secretCipher,
		Enabled:       false,
		RecoveryCodes: string(recoveryJSON),
	}
	if err := s.repo.Upsert(ctx, record); err != nil {
		return nil, err
	}
	return &TOTPBindResult{
		Secret:        key.Secret(),
		OTPAuthURL:    key.URL(),
		RecoveryCodes: recoveryCodes,
	}, nil
}

// ConfirmEnrollPending：用同一 step_token + 6 位码激活绑定并签发 access_token。
func (s *totpService) ConfirmEnrollPending(ctx context.Context, stepToken, code string) (*LoginResult, error) {
	claims, err := s.manager.VerifyStepToken(stepToken, core.StepTokenPurposeTOTPBind)
	if err != nil {
		return nil, errs.Unauthorized("step token 无效或已过期")
	}
	if used, _ := s.isStepJTIUsed(ctx, claims.ID); used {
		return nil, errs.Unauthorized("step token 已被使用")
	}

	record, err := s.repo.GetByAdminID(ctx, claims.AdminID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, errs.TOTPNotBound("尚未发起绑定，请先调用 enroll")
	}
	if record.Enabled {
		return nil, errs.TOTPAlreadyBound("已启用，无需重复确认")
	}
	secret, err := s.cipher.DecryptFromString(record.SecretCipher)
	if err != nil {
		return nil, errs.Internal("解密 TOTP secret 失败", err)
	}
	if !s.validateCode(secret, code) {
		return nil, errs.TOTPInvalid("验证码错误")
	}
	if err := s.repo.MarkEnabled(ctx, claims.AdminID); err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	_ = s.repo.UpdateLastUsedAt(ctx, claims.AdminID, now)
	s.markStepJTIUsed(ctx, claims.ID, time.Until(claims.ExpiresAt.Time))

	admin, err := s.adminRepo.GetAdminByID(ctx, claims.AdminID)
	if err != nil {
		return nil, errs.Unauthorized("用户不存在")
	}
	subject := core.Subject{UserID: admin.ID, Username: admin.Username}
	token, expires, err := s.manager.SignAccessToken(s.audience, subject)
	if err != nil {
		return nil, errs.Internal("生成 token 失败", err)
	}
	return &LoginResult{
		Stage:   LoginStageOK,
		Token:   token,
		Expires: expires,
		User:    &admin,
	}, nil
}

