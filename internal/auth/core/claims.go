package core

// 本文件定义 JWT subject、claims 与 audience 名称常量。

import "github.com/golang-jwt/jwt/v5"

// Audience 名称常量。这些是查找 Manager.audiences 的键，
// 与 JWT `aud` claim 的实际字符串通过配置映射而来。
const (
	// AudienceUser 表示 ToC API 终端用户。
	AudienceUser = "user"
	// AudienceAdmin 表示后台管理员。
	AudienceAdmin = "admin"
)

// Subject 表示需要为其签发 token 的身份。
type Subject struct {
	UserID    int64
	Username  string
	ProjectID string
}

// Claims 描述 access token 中保存的用户身份声明。
type Claims struct {
	UserID    int64  `json:"uid"`
	Username  string `json:"username"`
	ProjectID string `json:"project_id"`
	jwt.RegisteredClaims
}

// HasAudience 判断 claims 是否包含指定 audience 字符串。
func (c *Claims) HasAudience(audience string) bool {
	for _, a := range c.Audience {
		if a == audience {
			return true
		}
	}
	return false
}
