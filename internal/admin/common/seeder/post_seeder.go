package seeder

import (
	"gotribe/internal/model"

	"gorm.io/gorm"
)

// PostSeeder 文章种子
type PostSeeder struct {
	*BaseSeeder
}

// NewPostSeeder 创建文章种子
func NewPostSeeder() *PostSeeder {
	return &PostSeeder{
		BaseSeeder: NewBaseSeeder("post"),
	}
}

// Run 执行文章数据种子
func (s *PostSeeder) Run(db *gorm.DB, syncExisting bool) error {
	posts := []*model.Post{
		{
			Model:       model.Model{ID: 1},
			Slug:        "welcome-gotribe",
			Title:       "欢迎使用GoTribe",
			Description: "这是一篇示例文章",
			Content:     "# 这是一篇示例文章",
			Icon:        "https://cdn.example.com/sample.jpg",
			HtmlContent: "<h1>这是一篇示例文章</h1>",
			UserID:      1,
			CategoryID:  1,
			Author:      "GoTribe",
			ProjectID:   1,
		},
	}

	for _, post := range posts {
		if err := createIfNotExists(db, post, post.ID); err != nil {
			return err
		}
	}

	return nil
}
