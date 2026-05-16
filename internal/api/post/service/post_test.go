package service

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gotribe/internal/core/cache"
	"gotribe/internal/core/database"
	postmodel "gotribe/internal/model"
	postrepo "gotribe/internal/api/post/repository"
	tagrepo "gotribe/internal/api/tag/repository"
)

func TestService_DetailRejectsWrongPassword(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&postmodel.Post{}))
	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS post_tag (post_id INTEGER, tag_id INTEGER)`).Error)

	srv := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: srv.Addr()})
	defer redisClient.Close()

	tx := database.NewTransactionManager(db)
	keys := cache.NewKeyBuilder("gotribe-test")
	store := cache.NewStore(redisClient, keys)
	service := NewService(postrepo.NewRepository(tx), tagrepo.NewRepository(tx), store, 5)

	post := postmodel.Post{
		Slug:      "slug-protected",
		PostID:    "post-1",
		ProjectID: 1,
		Title:     "Protected",
		Status:    2,
		IsPasswd:  1,
		PassWord:  "secret",
	}
	require.NoError(t, db.Create(&post).Error)

	_, err = service.Detail(context.Background(), "1", "post-1", "wrong")
	require.Error(t, err)
}

func TestService_DetailCachesPublicPost(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&postmodel.Post{}))
	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS post_tag (post_id INTEGER, tag_id INTEGER)`).Error)

	srv := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: srv.Addr()})
	defer redisClient.Close()

	tx := database.NewTransactionManager(db)
	keys := cache.NewKeyBuilder("gotribe-test")
	store := cache.NewStore(redisClient, keys)
	service := NewService(postrepo.NewRepository(tx), tagrepo.NewRepository(tx), store, 5)

	post := postmodel.Post{
		Slug:      "slug-public",
		PostID:    "post-2",
		ProjectID: 1,
		Title:     "Public",
		Content:   "body",
		Status:    2,
	}
	require.NoError(t, db.Create(&post).Error)

	resp, err := service.Detail(context.Background(), "1", "post-2", "")
	require.NoError(t, err)
	require.Equal(t, "Public", resp.Title)

	// Content should be included for posts with no password
	require.Equal(t, "body", resp.Content)
}
