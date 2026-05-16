// Package view defines reusable internal view structures for the post module.
package view

// 本文件定义 post 模块内部对外暴露的文章摘要视图。

// Summary 表示供其他模块复用的文章摘要信息。
type Summary struct {
	PostID string
	Title  string
	Type   int16
	Status int16
}
