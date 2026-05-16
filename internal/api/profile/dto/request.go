// Package dto defines request and response structures for the profile module.
package dto

// 本文件定义个人资料模块的请求结构（修改密码、更新资料）。

// ChangePasswordRequest 表示修改当前用户密码接口的请求参数。
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// UpdateProfileRequest 表示更新当前用户资料接口的请求参数。
type UpdateProfileRequest struct {
	Nickname   *string `json:"nickname"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Sex        *string `json:"sex"`
	Birthday   *string `json:"birthday"`
	Background *string `json:"background"`
	Ext        *string `json:"ext"`
	AvatarURL  *string `json:"avatar_url"`
}
