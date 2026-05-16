package dto

// JobVO 任务视图对象
type JobVO struct {
	Name        string `json:"name"`        // 任务名称
	Description string `json:"description"` // 任务描述
	Schedule    string `json:"schedule"`    // 调度表达式
	Enabled     bool   `json:"enabled"`     // 是否启用
	Timeout     string `json:"timeout"`     // 超时时间
	RetryCount  int    `json:"retry_count"` // 重试次数
}
