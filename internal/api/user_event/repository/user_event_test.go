package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gotribe/internal/core/database"
	usereventmodel "gotribe/internal/model"
)

func TestRepository_CreatePersistsEvent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&usereventmodel.UserEvent{}))

	repo := NewRepository(database.NewTransactionManager(db))

	event := &usereventmodel.UserEvent{
		UserID:      42,
		ProjectID:   0,
		EventType:   2,
		EventDetail: "clicked button",
		Duration:    15,
		IP:          "127.0.0.1",
		UserAgent:   "test-agent",
		Referer:     "https://example.com",
		Platform:    "web",
	}
	require.NoError(t, repo.Create(context.Background(), event))

	var persisted usereventmodel.UserEvent
	require.NoError(t, db.First(&persisted, event.ID).Error)
	require.Equal(t, int64(0), persisted.ProjectID)
	require.EqualValues(t, 2, persisted.EventType)
	require.Equal(t, "clicked button", persisted.EventDetail)
}
