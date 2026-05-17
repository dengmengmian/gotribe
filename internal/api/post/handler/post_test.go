package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gotribe/internal/core/cache"
	"gotribe/internal/core/database"
	"gotribe/internal/core/middleware"
	postmodel "gotribe/internal/model"
	postrepo "gotribe/internal/api/post/repository"
	postservice "gotribe/internal/api/post/service"
	tagrepo "gotribe/internal/api/tag/repository"
	categoryrepo "gotribe/internal/api/category/repository"
)

func TestHandler_Detail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&postmodel.Post{}))
	require.NoError(t, db.Exec(`CREATE TABLE IF NOT EXISTS post_tag (post_id INTEGER, tag_id INTEGER)`).Error)

	post := postmodel.Post{
		Slug:      "slug-test",
		PostID:    "post-1",
		ProjectID: 1,
		Title:     "Hello",
		Content:   "body",
		Status:    2,
	}
	require.NoError(t, db.Create(&post).Error)

	srv := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: srv.Addr()})
	defer redisClient.Close()

	tx := database.NewTransactionManager(db)
	keys := cache.NewKeyBuilder("gotribe-test")
	store := cache.NewStore(redisClient, keys)
	handler := NewHandler(postservice.NewService(postrepo.NewRepository(tx), tagrepo.NewRepository(tx), categoryrepo.NewRepository(tx), store, 5))

	engine := gin.New()
	engine.Use(middleware.ProjectID("1"))
	engine.GET("/posts/:postID", handler.Detail)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/posts/post-1", nil)
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "Hello")
}
