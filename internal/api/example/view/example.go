// Package view defines reusable internal view structures for the example module.
package view

// 本文件定义 example 模块内部复用的视图结构。

// Actor 表示当前请求上下文中的操作者信息。
type Actor struct {
	UserID   int64
	Username string
	Nickname string
}

// Owner 表示示例业务单归属用户的快照信息。
type Owner struct {
	UserID   int64
	Username string
	Nickname string
}

// PostRef 表示示例业务单关联的文章摘要。
type PostRef struct {
	PostID string
	Title  string
	Type   int16
	Status int16
}

// Example 表示 example 模块服务层向外暴露的完整视图。
type Example struct {
	ExampleID   string
	Name        string
	Description string
	Status      int16
	Owner       Owner
	PrimaryPost PostRef
	Posts       []PostRef
	CreatedAt   string
	UpdatedAt   string
}
