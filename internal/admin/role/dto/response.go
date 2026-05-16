package dto

import "gotribe/internal/model"

// RoleResponse 角色列表/详情响应。
type RoleResponse struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Keyword   string `json:"keyword"`
	Desc      string `json:"desc"`
	Status    uint   `json:"status"`
	Sort      uint   `json:"sort"`
	Creator   string `json:"creator"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ToRoleResponse 将 model.Role 转为响应结构。
func ToRoleResponse(role *model.Role) RoleResponse {
	desc := ""
	if role.Desc != nil {
		desc = *role.Desc
	}
	return RoleResponse{
		ID:        role.ID,
		Name:      role.Name,
		Keyword:   role.Keyword,
		Desc:      desc,
		Status:    role.Status,
		Sort:      role.Sort,
		Creator:   role.Creator,
		CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToRoleListResponse 批量转换。
func ToRoleListResponse(roles []model.Role) []RoleResponse {
	res := make([]RoleResponse, 0, len(roles))
	for i := range roles {
		res = append(res, ToRoleResponse(&roles[i]))
	}
	return res
}

// MenuSummary 角色关联菜单摘要。
type MenuSummary struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Title string `json:"title"`
	Path  string `json:"path"`
}

// ToMenuSummaryList 批量转换菜单。
func ToMenuSummaryList(menus []*model.Menu) []MenuSummary {
	res := make([]MenuSummary, 0, len(menus))
	for _, m := range menus {
		if m == nil {
			continue
		}
		res = append(res, MenuSummary{
			ID:    m.ID,
			Name:  m.Name,
			Title: m.Title,
			Path:  m.Path,
		})
	}
	return res
}

// ApiSummary 角色关联 API 摘要。
type ApiSummary struct {
	ID       int64  `json:"id"`
	Method   string `json:"method"`
	Path     string `json:"path"`
	Category string `json:"category"`
	Desc     string `json:"desc"`
}

// ToApiSummaryList 批量转换 API。
func ToApiSummaryList(apis []*model.Api) []ApiSummary {
	res := make([]ApiSummary, 0, len(apis))
	for _, a := range apis {
		if a == nil {
			continue
		}
		res = append(res, ApiSummary{
			ID:       a.ID,
			Method:   a.Method,
			Path:     a.Path,
			Category: a.Category,
			Desc:     a.Desc,
		})
	}
	return res
}
