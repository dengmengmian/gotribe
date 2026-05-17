package core

// 本文件实现 TOTP 二次校验所需的短期 step_token：
//   - 是一个独立 JWT，purpose=totp_verify
//   - jti 一次性使用，外部用 Redis 黑名单防重放
//   - 不复用 access token 的 audience，避免被当作普通登录态使用

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// StepTokenPurposeTOTPVerify 是 TOTP 二次校验场景的 purpose 标记。
const StepTokenPurposeTOTPVerify = "totp_verify"

// stepTokenAudience 是 step_token 专属 audience，与登录 token 的 audience 隔离。
const stepTokenAudience = "gotribe.admin.step"

// StepClaims 描述 step_token 的 claims。
type StepClaims struct {
	AdminID  int64  `json:"aid"`
	Username string `json:"username"`
	Purpose  string `json:"purpose"`
	jwt.RegisteredClaims
}

// SignStepToken 为指定 admin 签发 step_token，TTL 由调用方传入（建议 5 分钟）。
// 返回 token、jti、过期时间。jti 调用方应写入 Redis 黑名单作单次使用控制。
func (m *Manager) SignStepToken(adminID int64, username, purpose string, ttl time.Duration) (string, string, time.Time, error) {
	if ttl <= 0 {
		return "", "", time.Time{}, fmt.Errorf("step token ttl must be positive")
	}
	jti, err := newJTI()
	if err != nil {
		return "", "", time.Time{}, err
	}
	expiresAt := time.Now().Add(ttl)
	claims := StepClaims{
		AdminID:  adminID,
		Username: username,
		Purpose:  purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   username,
			Audience:  jwt.ClaimStrings{stepTokenAudience},
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return signed, jti, expiresAt, nil
}

// VerifyStepToken 校验 step_token：issuer/audience/expiry/purpose 必须全部匹配。
// 调用方仍需检查 jti 是否已在 Redis 黑名单中（防重放）。
func (m *Manager) VerifyStepToken(tokenString, expectedPurpose string) (*StepClaims, error) {
	parsed, err := jwt.ParseWithClaims(
		tokenString,
		&StepClaims{},
		func(t *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(stepTokenAudience),
	)
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*StepClaims)
	if !ok || !parsed.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if claims.Purpose != expectedPurpose {
		return nil, errors.New("step token purpose mismatch")
	}
	return claims, nil
}

// newJTI 生成 128-bit 随机 ID（base64 url 编码）。
func newJTI() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
