package dto

import "gotribe/internal/model"

// AdminInfoResponse 返回给前端的当前用户信息
type AdminInfoResponse struct {
	ID           int64          `json:"id"`
	Username     string         `json:"username"`
	Mobile       string         `json:"mobile"`
	Avatar       string         `json:"avatar"`
	Nickname     string         `json:"nickname"`
	Introduction string         `json:"introduction"`
	Status       uint           `json:"status"`
	Roles        []*model.Role  `json:"roles"`
}

// ToAdminInfoResponse 转换为AdminInfoResponse
func ToAdminInfoResponse(user model.Admin) AdminInfoResponse {
	return AdminInfoResponse{
		ID:           user.ID,
		Username:     user.Username,
		Mobile:       user.Mobile,
		Avatar:       user.Avatar,
		Nickname:     *user.Nickname,
		Introduction: *user.Introduction,
		Status:       user.Status,
		Roles:        user.Roles,
	}
}

// AdminListResponse 返回给前端的用户列表
type AdminListResponse struct {
	ID           int64   `json:"id"`
	Username     string  `json:"username"`
	Mobile       string  `json:"mobile"`
	Avatar       string  `json:"avatar"`
	Nickname     string  `json:"nickname"`
	Introduction string  `json:"introduction"`
	Status       uint    `json:"status"`
	Creator      string  `json:"creator"`
	RoleIds      []int64 `json:"role_ids"`
}

// ToAdminListResponse 转换为AdminListResponse列表
func ToAdminListResponse(userList []*model.Admin) []AdminListResponse {
	var users []AdminListResponse
	for _, user := range userList {
		userDto := AdminListResponse{
			ID:           user.ID,
			Username:     user.Username,
			Mobile:       user.Mobile,
			Avatar:       user.Avatar,
			Nickname:     *user.Nickname,
			Introduction: *user.Introduction,
			Status:       user.Status,
			Creator:      user.Creator,
		}
		roleIds := make([]int64, 0)
		for _, role := range user.Roles {
			roleIds = append(roleIds, role.ID)
		}
		userDto.RoleIds = roleIds
		users = append(users, userDto)
	}

	return users
}
