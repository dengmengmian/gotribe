package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gotribe/internal/core/database"
	postmodel "gotribe/internal/model"
)

func newTestPostRepository(t *testing.T) (*Repository, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&postmodel.Post{}))
	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS post_tag (post_id INTEGER, tag_id INTEGER)`).Error)

	return NewRepository(database.NewTransactionManager(db)), db
}

func TestRepository_ListUsesDefaultPublishedScopeAndTagFilter(t *testing.T) {
	repo, db := newTestPostRepository(t)

	posts := []postmodel.Post{
		{Slug: "slug-pub-tagged", PostID: "post-published-tagged", ProjectID: 1, Title: "Published tagged", Status: 2, Type: 1, Sort: 20},
		{Slug: "slug-draft-tagged", PostID: "post-draft-tagged", ProjectID: 1, Title: "Draft tagged", Status: 1, Type: 1, Sort: 10},
		{Slug: "slug-pub-other", PostID: "post-published-other-project", ProjectID: 2, Title: "Other project", Status: 2, Type: 1, Sort: 30},
	}
	for i := range posts {
		require.NoError(t, db.Create(&posts[i]).Error)
	}

	require.NoError(t, db.Exec(`INSERT INTO post_tag (post_id, tag_id) VALUES (?, ?), (?, ?)`,
		posts[0].ID, 7, posts[1].ID, 7,
	).Error)

	got, total, err := repo.List(context.Background(), "1", ListFilter{
		Page:    1,
		PerPage: 10,
		TagIDs:  []int64{7},
	})
	require.NoError(t, err)
	require.EqualValues(t, 1, total)
	require.Len(t, got, 1)
	require.Equal(t, "post-published-tagged", got[0].PostID)
}

func TestRepository_GetByPostIDAndListByPostIDsRespectProjectScope(t *testing.T) {
	repo, db := newTestPostRepository(t)

	first := postmodel.Post{Slug: "slug-one", PostID: "post-1", ProjectID: 1, Title: "One", Status: 2, DynamicType: "news"}
	second := postmodel.Post{Slug: "slug-two", PostID: "post-2", ProjectID: 1, Title: "Two", Status: 2, DynamicType: "news"}
	other := postmodel.Post{Slug: "slug-three", PostID: "post-3", ProjectID: 2, Title: "Other", Status: 2, DynamicType: "news"}
	require.NoError(t, db.Create(&first).Error)
	require.NoError(t, db.Create(&second).Error)
	require.NoError(t, db.Create(&other).Error)

	got, err := repo.GetByPostID(context.Background(), "1", "post-1")
	require.NoError(t, err)
	require.Equal(t, "One", got.Title)

	_, err = repo.GetByPostID(context.Background(), "1", "post-3")
	require.Error(t, err)

	list, err := repo.ListByPostIDs(context.Background(), "1", []string{"post-1", "post-2", "post-3"})
	require.NoError(t, err)
	require.Len(t, list, 2)
	for _, item := range list {
		require.Equal(t, int64(1), item.ProjectID)
	}
}
