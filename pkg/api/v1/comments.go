// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package v1

// CreateCommentRequest 指定了 `POST /v1/comment` 接口的请求参数.
type CreateCommentRequest struct {
	Content    string `json:"content" valid:"required,stringlength(1|10240)"`
	ObjectID   string `json:"objectID" valid:"required,stringlength(1|10)"`
	ObjectType *uint  `json:"objectType" valid:"required"`
}

type ReplyCommentRequest struct {
	Content    string `json:"content" valid:"required,stringlength(1|10240)"`
	ObjectID   string `json:"objectID" valid:"required,stringlength(1|10)"`
	ObjectType uint   `json:"objectType"`
	ToUserID   string `json:"toUserID" valid:"required,stringlength(1|10)"`
	ParentID   *int   `json:"parentID" valid:"required"`
	ReplyToID  *int   `json:"replyToID""`
}

// CreateCommentResponse 指定了 `POST /v1/comment` 接口的返回参数.
type CreateCommentResponse struct {
	CommentID string `json:"commentID"`
}

// GetCommentResponse 指定了 `GET /v1/comment/{commentID}` 接口的返回参数.
type GetCommentResponse CommentInfo

// UpdateCommentRequest 指定了 `PUT /v1/comment` 接口的请求参数.
type UpdateCommentRequest struct {
	Content *string `json:"content" valid:"required,stringlength(1|10240)"`
}

// CommentInfo 指定了文章的详细信息.
type CommentInfo struct {
	ID          int            `json:"id"`
	CommentID   string         `json:"commentID"`
	Content     string         `json:"content" `
	HtmlContent string         `json:"htmlContent"`
	ObjectID    string         `json:"objectID"`
	ObjectType  uint           `json:"objectType"`
	ToUserID    string         `json:"toUserID" `
	ParentID    int            `json:"parent_id" `
	ReplyToID   int            `json:"replyToID"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
	UserID      string         `json:"user_id"`
	Nickname    string         `json:"nickname"`
	Avatar      string         `json:"avatar"`
	Replies     []*CommentInfo `json:"replies"`
	Country     string         `json:"country"`
	RegionName  string         `json:"regionName"`
	City        string         `json:"city"`
}

// ListCommentRequest 指定了 `GET /v1/comment` 接口的请求参数.
type ListCommentRequest struct {
	ObjectID string `form:"objectID" valid:"required,stringlength(1|10)"`
	Offset   int    `form:"offset"`
	Limit    int    `form:"limit"`
}

// ListCommentResponse 指定了 `GET /v1/comment` 接口的返回参数.
type ListCommentResponse struct {
	TotalCount int64          `json:"totalCount"`
	Comments   []*CommentInfo `json:"comment"`
}
