package dto

// GenerateRequest 是通用 AI 生成入口请求。
type GenerateRequest struct {
	Task     string         `json:"task" binding:"required"`
	Language string         `json:"language"`
	Input    map[string]any `json:"input" binding:"required"`
}
