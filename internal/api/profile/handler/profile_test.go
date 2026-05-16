package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gotribe/internal/auth/core"
	"gotribe/internal/core/cache"
	"gotribe/internal/core/database"
	"gotribe/internal/core/middleware"
	profilemodel "gotribe/internal/model"
	profilerepo "gotribe/internal/api/profile/repository"
	profileservice "gotribe/internal/api/profile/service"
	usermodel "gotribe/internal/model"
)

func TestHandler_GetMe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&profilemodel.UserProfile{}))

	user := profilemodel.UserProfile{
		Core: usermodel.Core{
			Username:  "tester",
			ProjectID: "proj-1",
			Password:  "hashed",
			Nickname:  "Tester",
			Status:    1,
		},
	}
	require.NoError(t, db.Create(&user).Error)

	srv := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: srv.Addr()})
	defer redisClient.Close()

	tx := database.NewTransactionManager(db)
	keys := cache.NewKeyBuilder("gotribe-test")
	store := cache.NewStore(redisClient, keys)
	tokens := core.NewTokenStore(redisClient, keys)
	handler := NewHandler(profileservice.NewService(core.AudienceUser, time.Hour, profilerepo.NewRepository(tx), store, tokens, tx, 5))

	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		c.Set(middleware.ContextKeyAuth, &middleware.AuthContext{UserID: int64(user.ID), ProjectID: "proj-1", Username: "tester"})
		c.Next()
	})
	engine.GET("/me", handler.GetMe)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	engine.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "tester")
}
