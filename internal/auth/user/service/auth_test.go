package service

// 本文件验证认证服务中 refresh 相关的安全边界。

import (
	"context"
	"errors"
	"testing"
	"time"

	"gotribe/internal/auth/core"
	"gotribe/internal/auth/user/dto"
	"gotribe/internal/core/errs"
	"gotribe/internal/model"
)

type fakeUserReader struct {
	findByIdentityFunc func(context.Context, string, string) (*model.AuthUser, error)
	findByIDFunc       func(context.Context, string, int64) (*model.AuthUser, error)
}

func (f fakeUserReader) FindByIdentity(ctx context.Context, projectID, identity string) (*model.AuthUser, error) {
	if f.findByIdentityFunc == nil {
		return nil, errors.New("unexpected FindByIdentity call")
	}
	return f.findByIdentityFunc(ctx, projectID, identity)
}

func (f fakeUserReader) FindByID(ctx context.Context, projectID string, userID int64) (*model.AuthUser, error) {
	if f.findByIDFunc == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return f.findByIDFunc(ctx, projectID, userID)
}

type fakeRefreshTokenStore struct {
	session    core.RefreshSession
	ok         bool
	getErr     error
	rotateOK   bool
	rotateErr  error
	rotateCall int
}

func (f *fakeRefreshTokenStore) Save(context.Context, string, string, core.RefreshSession, time.Duration) error {
	return errors.New("unexpected Save call")
}

func (f *fakeRefreshTokenStore) Get(context.Context, string, string) (core.RefreshSession, bool, error) {
	return f.session, f.ok, f.getErr
}

func (f *fakeRefreshTokenStore) Delete(context.Context, string, string) error {
	return errors.New("unexpected Delete call")
}

func (f *fakeRefreshTokenStore) Rotate(context.Context, string, string, core.RefreshSession, string, core.RefreshSession, time.Duration) (bool, error) {
	f.rotateCall++
	return f.rotateOK, f.rotateErr
}

type fakeTokenManager struct{}

func (fakeTokenManager) SignAccessToken(string, core.Subject) (string, time.Time, error) {
	return "access", time.Now(), nil
}

func (fakeTokenManager) GenerateRefreshToken() (string, error) {
	return "refresh", nil
}

func (fakeTokenManager) AccessTTL(string) (time.Duration, bool) {
	return time.Hour, true
}

func (fakeTokenManager) RefreshTTL(string) (time.Duration, bool) {
	return 24 * time.Hour, true
}

// TestRefreshRejectsDisabledUser 验证禁用用户不能通过 refresh token 继续换发会话。
func TestRefreshRejectsDisabledUser(t *testing.T) {
	t.Parallel()

	store := &fakeRefreshTokenStore{
		session: core.RefreshSession{Audience: core.AudienceUser, UserID: 1, ProjectID: "demo"},
		ok:      true,
	}
	svc := NewService(
		core.AudienceUser,
		fakeUserReader{
			findByIDFunc: func(context.Context, string, int64) (*model.AuthUser, error) {
				return &model.AuthUser{Core: model.Core{Username: "tester", ProjectID: "demo", Status: 0}}, nil
			},
		},
		store,
		fakeTokenManager{},
	)

	_, err := svc.Refresh(context.Background(), dto.RefreshRequest{RefreshToken: "valid"})
	appErr := errs.As(err)
	if appErr == nil || appErr.Code != errs.CodeForbidden {
		t.Fatalf("error = %#v, want forbidden app error", appErr)
	}
	if store.rotateCall != 0 {
		t.Fatalf("Rotate() called %d times, want 0", store.rotateCall)
	}
}

// TestRefreshKeepsTokenWhenUserLookupFails 验证用户查询失败时不会提前旋转掉旧 refresh token。
func TestRefreshKeepsTokenWhenUserLookupFails(t *testing.T) {
	t.Parallel()

	store := &fakeRefreshTokenStore{
		session: core.RefreshSession{Audience: core.AudienceUser, UserID: 1, ProjectID: "demo"},
		ok:      true,
	}
	svc := NewService(
		core.AudienceUser,
		fakeUserReader{
			findByIDFunc: func(context.Context, string, int64) (*model.AuthUser, error) {
				return nil, errors.New("db unavailable")
			},
		},
		store,
		fakeTokenManager{},
	)

	_, err := svc.Refresh(context.Background(), dto.RefreshRequest{RefreshToken: "valid"})
	appErr := errs.As(err)
	if appErr == nil || appErr.Code != errs.CodeInternal {
		t.Fatalf("error = %#v, want internal app error", appErr)
	}
	if store.rotateCall != 0 {
		t.Fatalf("Rotate() called %d times, want 0", store.rotateCall)
	}
}
