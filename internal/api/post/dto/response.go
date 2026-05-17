package dto

// 本文件定义帖子模块的响应结构。

// TagResponse 表示文章接口返回的标签信息。
type TagResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
	Color string `json:"color"`
}

// CategoryResponse 表示文章接口返回的分类信息。
type CategoryResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// PostResponse 表示文章列表或详情接口的响应数据。
type PostResponse struct {
	ID           int64         `json:"id"`
	PostID       string        `json:"post_id"`
	Slug         string        `json:"slug"`
	CategoryID   int64         `json:"category_id"`
	ProjectID    string        `json:"project_id"`
	UserID       int64         `json:"user_id"`
	Author       string        `json:"author"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	HTMLContent  string        `json:"html_content"`
	WordCount    int           `json:"word_count"`
	Description  string        `json:"description"`
	Icon         string        `json:"icon"`
	View         int64         `json:"view"`
	Type         int16         `json:"type"`
	Status       int16         `json:"status"`
	UnitPrice    int           `json:"unit_price"`
	Location     string        `json:"location"`
	People       string        `json:"people"`
	Time         string        `json:"time"`
	ShowTime     string        `json:"show_time"`
	DynamicType  string        `json:"dynamic_type"`
	Sort         uint          `json:"sort"`
	EventStartAt string        `json:"event_start_at"`
	EventEndAt   string        `json:"event_end_at"`
	RegisterURL  string        `json:"register_url"`
	Tags         []TagResponse     `json:"tags"`
	Category     *CategoryResponse `json:"category,omitempty"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
}
