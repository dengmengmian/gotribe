package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"gotribe/internal/auth/core"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	"gotribe/internal/model"
	utils "gotribe/internal/core/util"
)

func setupTestDB(t *testing.T) (*gorm.DB, *database.TransactionManager) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&model.Admin{}, &model.Role{})
	require.NoError(t, err)
	return db, database.NewTransactionManager(db)
}

func newTestManager(t *testing.T) *core.Manager {
	t.Helper()
	manager, err := core.NewManager("test", "test-secret-must-be-at-least-32-characters", map[string]core.AudienceConfig{
		core.AudienceAdmin: {
			Audience:        "test.admin",
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: 24 * time.Hour,
		},
	})
	require.NoError(t, err)
	return manager
}

func createTestAdmin(t *testing.T, db *gorm.DB) *model.Admin {
	hashedPassword, err := utils.PasswordUtil.GenPasswd("Gotribe!23456")
	require.NoError(t, err)

	role := &model.Role{
		Name:    "TestRole",
		Keyword: "test_role",
		Status:  1,
		Sort:    100,
	}
	err = db.Create(role).Error
	require.NoError(t, err)

	admin := &model.Admin{
		Username: "testadmin",
		Password: hashedPassword,
		Status:   1,
		Roles:    []*model.Role{role},
	}
	err = db.Create(admin).Error
	require.NoError(t, err)
	return admin
}

func TestNewService(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	tx := database.NewTransactionManager(db)
	svc := NewService(core.AudienceAdmin, tx, newTestManager(t), nil, nil, 5*time.Minute, false)
	require.NotNil(t, svc)
}

func TestAuthService_Login(t *testing.T) {
	db, tx := setupTestDB(t)
	admin := createTestAdmin(t, db)

	manager := newTestManager(t)
	svc := NewService(core.AudienceAdmin, tx, manager, nil, nil, 5*time.Minute, false)
	ctx := context.Background()

	t.Run("登录成功", func(t *testing.T) {
		result, err := svc.Login(ctx, admin.Username, "Gotribe!23456", "127.0.0.1")
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotEmpty(t, result.Token)
		require.NotZero(t, result.Expires)
		require.Equal(t, admin.Username, result.User.Username)

		// 验证生成的 token 可以被解析
		claims, err := manager.VerifyAccessToken(core.AudienceAdmin, result.Token)
		require.NoError(t, err)
		require.Equal(t, admin.Username, claims.Username)
	})

	t.Run("密码错误", func(t *testing.T) {
		result, err := svc.Login(ctx, admin.Username, "wrongpassword", "127.0.0.1")
		require.Error(t, err)
		require.Nil(t, result)
		var ae *errs.AppError
		require.True(t, errors.As(err, &ae))
		require.Equal(t, errs.CodeUnauthorized, ae.Code)
	})

	t.Run("用户不存在", func(t *testing.T) {
		result, err := svc.Login(ctx, "nonexistent", "Gotribe!23456", "127.0.0.1")
		require.Error(t, err)
		require.Nil(t, result)
		var ae *errs.AppError
		require.True(t, errors.As(err, &ae))
		require.Equal(t, errs.CodeUnauthorized, ae.Code)
	})
}

func TestAuthService_Refresh(t *testing.T) {
	db, tx := setupTestDB(t)
	admin := createTestAdmin(t, db)

	manager := newTestManager(t)
	svc := NewService(core.AudienceAdmin, tx, manager, nil, nil, 5*time.Minute, false)
	ctx := context.Background()

	t.Run("刷新成功", func(t *testing.T) {
		result, err := svc.Refresh(ctx, admin.ID, admin.Username)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotEmpty(t, result.Token)
		require.NotZero(t, result.Expires)

		// 验证新 token 有效
		claims, err := manager.VerifyAccessToken(core.AudienceAdmin, result.Token)
		require.NoError(t, err)
		require.Equal(t, admin.Username, claims.Username)
	})

	t.Run("刷新生成不同 token", func(t *testing.T) {
		result1, err := svc.Refresh(ctx, admin.ID, admin.Username)
		require.NoError(t, err)

		time.Sleep(2 * time.Second) // 确保 iat 不同

		result2, err := svc.Refresh(ctx, admin.ID, admin.Username)
		require.NoError(t, err)

		require.NotEqual(t, result1.Token, result2.Token)
	})
}
