package core

// 本文件集中定义认证 core 包暴露的错误类型。

import "errors"

var (
	// ErrUnknownAudience 表示 Manager 中未配置该 audience。
	ErrUnknownAudience = errors.New("unknown audience")
	// ErrAudienceMismatch 表示 token 中 `aud` 与期望不匹配。
	ErrAudienceMismatch = errors.New("audience mismatch")
	// ErrIssuerMismatch 表示 token 中 `iss` 与期望不匹配。
	ErrIssuerMismatch = errors.New("issuer mismatch")
	// ErrInvalidBearerToken 表示 Authorization header 缺失或格式不合法。
	ErrInvalidBearerToken = errors.New("invalid bearer token")
)
