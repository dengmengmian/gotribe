// Package core provides a unified, audience-aware authentication core
// (JWT signing/verification, refresh token storage, password hashing,
// and gin middleware) shared by both the public API and Admin entry points.
package core

// 本文件实现 audience 参数化的 JWT 签发与校验。

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AudienceConfig 描述单个 audience 的 token 参数。
type AudienceConfig struct {
	// Audience 是写入 JWT `aud` claim 的字符串。
	Audience string
	// AccessTokenTTL 访问令牌有效期。
	AccessTokenTTL time.Duration
	// RefreshTokenTTL 刷新令牌有效期。
	RefreshTokenTTL time.Duration
}

// Manager 负责签发与校验跨 audience 的 JWT token。
// 同一 issuer 与 secret 下，多 audience 互相隔离（aud 不匹配即拒绝）。
type Manager struct {
	issuer    string
	secret    []byte
	audiences map[string]AudienceConfig
}

// NewManager 创建 Manager。
//   - issuer 必填
//   - secret 至少 32 字节
//   - audiences 至少一项，每项 Audience / AccessTokenTTL / RefreshTokenTTL 都需有效
func NewManager(issuer, secret string, audiences map[string]AudienceConfig) (*Manager, error) {
	if strings.TrimSpace(issuer) == "" {
		return nil, fmt.Errorf("issuer is required")
	}
	if len(strings.TrimSpace(secret)) < 32 {
		return nil, fmt.Errorf("secret must be at least 32 characters")
	}
	if len(audiences) == 0 {
		return nil, fmt.Errorf("at least one audience is required")
	}
	cloned := make(map[string]AudienceConfig, len(audiences))
	for name, cfg := range audiences {
		if strings.TrimSpace(cfg.Audience) == "" {
			return nil, fmt.Errorf("audience %q has empty Audience", name)
		}
		if cfg.AccessTokenTTL <= 0 {
			return nil, fmt.Errorf("audience %q access ttl must be positive", name)
		}
		if cfg.RefreshTokenTTL <= 0 {
			return nil, fmt.Errorf("audience %q refresh ttl must be positive", name)
		}
		cloned[name] = cfg
	}
	return &Manager{
		issuer:    issuer,
		secret:    []byte(secret),
		audiences: cloned,
	}, nil
}

// SignAccessToken 为指定 audience 签发 access token。
func (m *Manager) SignAccessToken(audienceName string, subject Subject) (string, time.Time, error) {
	cfg, ok := m.audiences[audienceName]
	if !ok {
		return "", time.Time{}, fmt.Errorf("%w: %q", ErrUnknownAudience, audienceName)
	}
	expiresAt := time.Now().Add(cfg.AccessTokenTTL)
	claims := Claims{
		UserID:    subject.UserID,
		Username:  subject.Username,
		ProjectID: subject.ProjectID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   subject.Username,
			Audience:  jwt.ClaimStrings{cfg.Audience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

// VerifyAccessToken 校验 access token；issuer 与 audience 必须匹配，过期会被拒绝。
func (m *Manager) VerifyAccessToken(audienceName, tokenString string) (*Claims, error) {
	cfg, ok := m.audiences[audienceName]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownAudience, audienceName)
	}
	parsed, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithIssuer(m.issuer),
		jwt.WithAudience(cfg.Audience),
	)
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

// VerifyAccessTokenWithoutExpiry 校验 token 但忽略过期；仍校验 issuer 与 audience。
// 用于 refresh 端点：在已过期的 access token 基础上读取身份信息。
func (m *Manager) VerifyAccessTokenWithoutExpiry(audienceName, tokenString string) (*Claims, error) {
	cfg, ok := m.audiences[audienceName]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownAudience, audienceName)
	}
	parsed, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithoutClaimsValidation(),
	)
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if claims.Issuer != m.issuer {
		return nil, ErrIssuerMismatch
	}
	if !claims.HasAudience(cfg.Audience) {
		return nil, ErrAudienceMismatch
	}
	return claims, nil
}

// GenerateRefreshToken 生成 256-bit 高强度随机 refresh token。
func (m *Manager) GenerateRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// AccessTTL 返回指定 audience 的 access token 时长。
func (m *Manager) AccessTTL(audienceName string) (time.Duration, bool) {
	cfg, ok := m.audiences[audienceName]
	if !ok {
		return 0, false
	}
	return cfg.AccessTokenTTL, true
}

// RefreshTTL 返回指定 audience 的 refresh token 时长。
func (m *Manager) RefreshTTL(audienceName string) (time.Duration, bool) {
	cfg, ok := m.audiences[audienceName]
	if !ok {
		return 0, false
	}
	return cfg.RefreshTokenTTL, true
}

// Audience 返回指定 audience 名称对应的实际 audience 字符串。
func (m *Manager) Audience(audienceName string) (string, bool) {
	cfg, ok := m.audiences[audienceName]
	return cfg.Audience, ok
}

// ParseBearerToken 从 Authorization header 中提取 Bearer 令牌字符串。
// 大小写不敏感（与原 core/middleware/jwt.go 行为保持一致）。
func ParseBearerToken(authHeader string) (string, error) {
	header := strings.TrimSpace(authHeader)
	if len(header) < 7 || !strings.EqualFold(header[:7], "Bearer ") {
		return "", ErrInvalidBearerToken
	}
	return strings.TrimSpace(header[7:]), nil
}
