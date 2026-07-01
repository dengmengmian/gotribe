package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	usereventdto "gotribe/internal/api/user_event/dto"
	usereventrepo "gotribe/internal/api/user_event/repository"
	"gotribe/internal/core/database"
	usereventmodel "gotribe/internal/model"
)

func TestService_CreatePersistsEvent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&usereventmodel.UserEvent{}))

	tx := database.NewTransactionManager(db)
	service := NewService(usereventrepo.NewRepository(tx))

	err = service.Create(context.Background(), "1", 1, usereventdto.CreateRequest{
		EventType:   1,
		EventDetail: "opened page",
		Duration:    10,
		Referer:     "https://example.com",
		Platform:    "web",
	}, "127.0.0.1", "ua")
	require.NoError(t, err)

	var count int64
	require.NoError(t, db.Model(&usereventmodel.UserEvent{}).Count(&count).Error)
	require.Equal(t, int64(1), count)
}
