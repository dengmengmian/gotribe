package dto

import (
	"gotribe/internal/model"
	"gotribe/internal/core/constant"
	"gotribe/internal/core/util"
	"strings"
	"time"
)

// PostResponse 返回给前端的内容
type PostResponse struct {
	ID          int64           `json:"id"`
	Slug        string          `json:"slug"`
	ColumnID    int64           `json:"column_id,omitempty"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	CategoryID  int64           `json:"category_id"`
	ProjectID   int64           `json:"project_id"`
	UserID      int64           `json:"user_id"`
	Author      string          `json:"author"`
	Content     string          `json:"content"`
	HtmlContent string          `json:"html_content"`
	Ext         string          `json:"ext"`
	Icon        string          `json:"icon"`
	Tag         string          `json:"tag"`
	Type        uint            `json:"type"`
	IsTop       int64           `json:"is_top"`
	IsPasswd    int64           `json:"is_passwd"`
	Category    *model.Category `json:"category"`
	Tags        []*model.Tag    `json:"tags"`
	Project     *model.Project  `json:"project"`
	CreatedAt   string          `json:"created_at"`
	Status      uint            `json:"status"`
	Location    string          `json:"location"`
	People      string          `json:"people"`
	Time        string          `json:"time"`
	Images      []string        `json:"images"`
	UnitPrice   float64         `json:"unit_price"`
	ShowTime    string          `json:"show_time"`
	Video       string          `json:"video"`
}

// ToPostResponse 将单个 Post 转换为 PostResponse
func ToPostResponse(post *model.Post) PostResponse {
	if post == nil {
		return PostResponse{}
	}
	var imageList []string
	if len(post.Images) > 0 {
		// 用,分割成数组
		imageList = strings.Split(post.Images, ",")
	}
	return PostResponse{
		ID:          post.ID,
		Slug:        post.Slug,
		ColumnID:    post.ColumnID,
		Title:       post.Title,
		Description: post.Description,
		CategoryID:  post.CategoryID,
		ProjectID:   post.ProjectID,
		UserID:      post.UserID,
		Author:      post.Author,
		Content:     post.Content,
		HtmlContent: post.HtmlContent,
		Ext:         post.Ext,
		Icon:        post.Icon,
		Tag:         post.Tag,
		Type:        post.Type,
		IsTop:       post.IsTop,
		IsPasswd:    post.IsPasswd,
		Category:    post.Category,
		CreatedAt:   post.CreatedAt.Format(constant.TIME_FORMAT),
		Tags:        post.Tags,
		Project:     post.Project,
		Status:      post.Status,
		Location:    post.Location,
		People:      post.People,
		Time:        formatPostTime(post.Time, constant.TIME_FORMAT_SHORT),
		Images:      imageList,
		UnitPrice:   utils.MoneyUtil.CentsToYuan(int64(post.UnitPrice)),
		Video:       post.Video,
		ShowTime:    formatPostTime(post.ShowTime, constant.TIME_FORMAT),
	}
}

// ToPostListResponse 将多个 Post 转换为 PostResponse 列表
func ToPostListResponse(postList []*model.Post) []PostResponse {
	var posts []PostResponse
	for _, post := range postList {
		if post == nil {
			continue
		}
		postResponse := ToPostResponse(post)
		posts = append(posts, postResponse)
	}

	return posts
}

func formatPostTime(value *time.Time, layout string) string {
	if value == nil {
		return ""
	}
	return value.Format(layout)
}
