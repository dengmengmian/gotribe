package service

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gotribe/internal/auth/core"
	"gotribe/internal/core/cache"
	"gotribe/internal/core/database"
	profiledto "gotribe/internal/api/profile/dto"
	profilemodel "gotribe/internal/model"
	profilerepo "gotribe/internal/api/profile/repository"
	usermodel "gotribe/internal/model"
)

func TestService_UpdateMePreservesSpecialCharacters(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&profilemodel.UserProfile{}))

	srv := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{Addr: srv.Addr()})
	defer redisClient.Close()

	tx := database.NewTransactionManager(db)
	repo := profilerepo.NewRepository(tx)
	keys := cache.NewKeyBuilder("gotribe-test")
	store := cache.NewStore(redisClient, keys)
	tokens := core.NewTokenStore(redisClient, keys)
	service := NewService(core.AudienceUser, time.Hour, repo, store, tokens, tx, 5)

	user := profilemodel.UserProfile{
		Core: usermodel.Core{
			Username:  "tester",
			ProjectID: "proj-1",
			Password:  "hashed",
			Nickname:  "old",
		},
	}
	require.NoError(t, db.Create(&user).Error)

	nickname := "Tom & Jerry"
	background := "A&B"
	ext := `{"a":"b&c"}`
	resp, err := service.UpdateMe(context.Background(), "proj-1", int64(user.ID), profiledto.UpdateProfileRequest{
		Nickname:   &nickname,
		Background: &background,
		Ext:        &ext,
	})
	require.NoError(t, err)
	require.Equal(t, nickname, resp.Nickname)
	require.Equal(t, background, resp.Background)
	require.Equal(t, ext, resp.Ext)
}
