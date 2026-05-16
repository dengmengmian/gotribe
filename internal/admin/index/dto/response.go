package dto

// IndexResponse 仪表盘全量数据。
type IndexResponse struct {
	Stats          Stats            `json:"stats"`
	VisitTrend     []VisitPoint     `json:"visit_trend"`
	RecentPosts    []PostSummary    `json:"recent_posts"`
	RecentComments []CommentSummary `json:"recent_comments"`
	PopularPosts   []PostSummary    `json:"popular_posts"`
	Pending        Pending          `json:"pending"`
	SystemStatus   SystemStatus     `json:"system_status"`
	CacheStatus    CacheStatus      `json:"cache_status"`
	SeoAlerts      []SeoAlert       `json:"seo_alerts"`
}

// Stats 统计卡片数据。
type Stats struct {
	TotalPosts      int64 `json:"total_posts"`
	DraftPosts      int64 `json:"draft_posts"`
	PendingComments int64 `json:"pending_comments"`
	WeekVisits      int64 `json:"week_visits"`
}

// VisitPoint 单日访问数据点。
type VisitPoint struct {
	Date      string `json:"date"`
	Visits    int64  `json:"visits"`
	PageViews int64  `json:"page_views"`
}

// PostSummary 文章摘要。
type PostSummary struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Status    uint   `json:"status"`
	CreatedAt string `json:"created_at"`
	View      int64  `json:"view"`
}

// CommentSummary 评论摘要。
type CommentSummary struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	Content   string `json:"content"`
	PostTitle string `json:"post_title"`
	CreatedAt string `json:"created_at"`
}

// Pending 待处理计数。
type Pending struct {
	ReviewPosts    int64 `json:"pending_review_posts"`
	ReviewComments int64 `json:"pending_review_comments"`
}

// SystemStatus 系统运行状态。
type SystemStatus struct {
	DBStatus    string `json:"db_status"`
	RedisStatus string `json:"redis_status"`
}

// CacheStatus Redis 缓存状态。
type CacheStatus struct {
	UsedMemory  string `json:"used_memory"`
	UsedPercent int    `json:"used_percent"`
}

// SeoAlert SEO 提醒。
type SeoAlert struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
