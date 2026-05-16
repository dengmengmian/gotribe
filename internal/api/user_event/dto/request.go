// Package dto defines request structures for the user event module.
package dto

// 本文件定义用户行为上报接口的请求结构。

// CreateRequest 表示用户行为事件上报接口的请求参数。
type CreateRequest struct {
	EventType   int16  `json:"event_type" binding:"required"`
	EventDetail string `json:"event_detail"`
	Duration    int    `json:"duration"`
	Referer     string `json:"referer"`
	Platform    string `json:"platform"`
}
