package dto

import (
	"gotribe/internal/core/constant"
	"gotribe/internal/model"
)

type MenuResponse struct {
	ID         int64          `json:"id"`
	Name       string         `json:"name"`
	Title      string         `json:"title"`
	Icon       *string        `json:"icon"`
	Path       string         `json:"path"`
	Redirect   *string        `json:"redirect"`
	Component  string         `json:"component"`
	Sort       uint           `json:"sort"`
	Status     uint           `json:"status"`
	Hidden     uint           `json:"hidden"`
	NoCache    uint           `json:"no_cache"`
	AlwaysShow uint           `json:"always_show"`
	Breadcrumb int64          `json:"breadcrumb"`
	ActiveMenu *string        `json:"active_menu"`
	ParentID   *int64         `json:"parent_id"`
	Creator    string         `json:"creator"`
	CreatedAt  string         `json:"created_at"`
	Children   []*MenuResponse `json:"children,omitempty"`
}

func toMenuResponse(menu *model.Menu) MenuResponse {
	if menu == nil {
		return MenuResponse{}
	}
	return MenuResponse{
		ID:         menu.ID,
		Name:       menu.Name,
		Title:      menu.Title,
		Icon:       menu.Icon,
		Path:       menu.Path,
		Redirect:   menu.Redirect,
		Component:  menu.Component,
		Sort:       menu.Sort,
		Status:     menu.Status,
		Hidden:     menu.Hidden,
		NoCache:    menu.NoCache,
		AlwaysShow: menu.AlwaysShow,
		Breadcrumb: menu.Breadcrumb,
		ActiveMenu: menu.ActiveMenu,
		ParentID:   menu.ParentID,
		Creator:    menu.Creator,
		CreatedAt:  menu.CreatedAt.Format(constant.TIME_FORMAT),
	}
}

func ToMenuResponse(menu model.Menu) MenuResponse {
	return toMenuResponse(&menu)
}

func ToMenuListResponse(menuList []*model.Menu) []MenuResponse {
	if menuList == nil {
		return []MenuResponse{}
	}

	menus := make([]MenuResponse, 0, len(menuList))
	for _, menu := range menuList {
		menus = append(menus, toMenuResponse(menu))
	}

	return menus
}

func toMenuTreeResponseList(menus []*model.Menu) []*MenuResponse {
	result := make([]*MenuResponse, 0, len(menus))
	for _, m := range menus {
		menuRes := toMenuResponse(m)
		if len(m.Children) > 0 {
			children := toMenuTreeResponseList(m.Children)
			menuRes.Children = children
		}
		result = append(result, &menuRes)
	}
	return result
}

func ToMenuTreeResponse(menus []*model.Menu) []*MenuResponse {
	return toMenuTreeResponseList(menus)
}
