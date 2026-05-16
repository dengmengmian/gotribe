// Package dto defines request and response structures for the authentication module.
package dto

// 本文件定义认证模块的请求结构（登录、登出、刷新、验证码）。

// LoginRequest 表示登录接口的请求参数。
type LoginRequest struct {
	Identity string `json:"identity" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LogoutRequest 表示退出登录接口的请求参数。
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshRequest 表示刷新令牌接口的请求参数。
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// VerifyCodeRequest 表示验证码相关接口的请求参数。
type VerifyCodeRequest struct {
	Target string `json:"target"`
	Scene  string `json:"scene"`
}
