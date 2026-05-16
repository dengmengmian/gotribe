package dto

// 本文件定义 example 模块的响应结构。

// OwnerResponse 表示示例业务单返回的归属用户摘要信息。
type OwnerResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

// PostRefResponse 表示示例业务单关联的文章摘要信息。
type PostRefResponse struct {
	PostID string `json:"post_id"`
	Title  string `json:"title"`
	Type   int16  `json:"type"`
	Status int16  `json:"status"`
}

// ExampleResponse 表示示例业务单接口的响应结构。
type ExampleResponse struct {
	ExampleID   string            `json:"example_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      int16             `json:"status"`
	Owner       OwnerResponse     `json:"owner"`
	PrimaryPost PostRefResponse   `json:"primary_post"`
	Posts       []PostRefResponse `json:"posts"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}
