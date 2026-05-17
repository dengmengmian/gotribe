package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func newTestManagerForStep(t *testing.T) *Manager {
	t.Helper()
	m, err := NewManager("test-issuer", "test-secret-must-be-at-least-32-characters", map[string]AudienceConfig{
		AudienceAdmin: {
			Audience:        "test.admin",
			AccessTokenTTL:  time.Hour,
			RefreshTokenTTL: 24 * time.Hour,
		},
	})
	require.NoError(t, err)
	return m
}

func TestStepToken_SignAndVerify(t *testing.T) {
	m := newTestManagerForStep(t)

	token, jti, expires, err := m.SignStepToken(42, "alice", StepTokenPurposeTOTPVerify, 5*time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, jti)
	require.WithinDuration(t, time.Now().Add(5*time.Minute), expires, 2*time.Second)

	claims, err := m.VerifyStepToken(token, StepTokenPurposeTOTPVerify)
	require.NoError(t, err)
	require.Equal(t, int64(42), claims.AdminID)
	require.Equal(t, "alice", claims.Username)
	require.Equal(t, jti, claims.ID)
}

func TestStepToken_RejectsWrongPurpose(t *testing.T) {
	m := newTestManagerForStep(t)
	token, _, _, err := m.SignStepToken(1, "alice", StepTokenPurposeTOTPVerify, 5*time.Minute)
	require.NoError(t, err)

	_, err = m.VerifyStepToken(token, "other_purpose")
	require.Error(t, err)
}

func TestStepToken_RejectsExpired(t *testing.T) {
	m := newTestManagerForStep(t)
	token, _, _, err := m.SignStepToken(1, "alice", StepTokenPurposeTOTPVerify, 1*time.Millisecond)
	require.NoError(t, err)
	time.Sleep(20 * time.Millisecond)

	_, err = m.VerifyStepToken(token, StepTokenPurposeTOTPVerify)
	require.Error(t, err)
}

func TestStepToken_RejectsZeroTTL(t *testing.T) {
	m := newTestManagerForStep(t)
	_, _, _, err := m.SignStepToken(1, "alice", StepTokenPurposeTOTPVerify, 0)
	require.Error(t, err)
}

func TestStepToken_JTIIsUnique(t *testing.T) {
	m := newTestManagerForStep(t)
	_, jti1, _, err := m.SignStepToken(1, "alice", StepTokenPurposeTOTPVerify, 5*time.Minute)
	require.NoError(t, err)
	_, jti2, _, err := m.SignStepToken(1, "alice", StepTokenPurposeTOTPVerify, 5*time.Minute)
	require.NoError(t, err)
	require.NotEqual(t, jti1, jti2)
}
