package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gotribe/internal/core/database"
	profilemodel "gotribe/internal/model"
	usermodel "gotribe/internal/model"
)

func newTestRepository(t *testing.T) (*Repository, *gorm.DB) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&profilemodel.UserProfile{}))

	return NewRepository(database.NewTransactionManager(db)), db
}

func TestRepository_GetByIDAndPasswordRespectProjectScope(t *testing.T) {
	repo, db := newTestRepository(t)

	user := profilemodel.UserProfile{
		Core: usermodel.Core{
			ProjectID: "proj-1",
			Username:  "tester",
			Password:  "hashed-password",
			Nickname:  "Tester",
		},
	}
	require.NoError(t, db.Create(&user).Error)

	got, err := repo.GetByID(context.Background(), "proj-1", int64(user.ID))
	require.NoError(t, err)
	require.Equal(t, user.ID, got.ID)
	require.Equal(t, "Tester", got.Nickname)

	password, err := repo.GetPasswordByID(context.Background(), "proj-1", int64(user.ID))
	require.NoError(t, err)
	require.Equal(t, "hashed-password", password)

	_, err = repo.GetByID(context.Background(), "proj-2", int64(user.ID))
	require.Error(t, err)

	_, err = repo.GetPasswordByID(context.Background(), "proj-2", int64(user.ID))
	require.Error(t, err)
}

func TestRepository_UpdateAndUpdatePasswordRespectProjectScope(t *testing.T) {
	repo, db := newTestRepository(t)

	user := profilemodel.UserProfile{
		Core: usermodel.Core{
			ProjectID: "proj-1",
			Username:  "tester",
			Password:  "old-password",
			Nickname:  "Old",
		},
		Background: "before",
	}
	require.NoError(t, db.Create(&user).Error)

	err := repo.Update(context.Background(), "proj-1", int64(user.ID), map[string]any{
		"nickname":   "New",
		"background": "after",
	})
	require.NoError(t, err)

	var updated profilemodel.UserProfile
	require.NoError(t, db.First(&updated, int64(user.ID)).Error)
	require.Equal(t, "New", updated.Nickname)
	require.Equal(t, "after", updated.Background)

	err = repo.UpdatePassword(context.Background(), "proj-2", int64(user.ID), "should-not-apply")
	require.NoError(t, err)
	require.NoError(t, db.First(&updated, int64(user.ID)).Error)
	require.Equal(t, "old-password", updated.Password)

	err = repo.UpdatePassword(context.Background(), "proj-1", int64(user.ID), "new-password")
	require.NoError(t, err)
	require.NoError(t, db.First(&updated, int64(user.ID)).Error)
	require.Equal(t, "new-password", updated.Password)
}
