package service

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gotribe/internal/core/database"
)

func TestNewUserService(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	svc := NewService(database.NewTransactionManager(db), nil, nil)
	require.NotNil(t, svc)
}
