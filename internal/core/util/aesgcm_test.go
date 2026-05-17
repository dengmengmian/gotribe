package utils

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAESGCMCipher_KeyLength(t *testing.T) {
	t.Run("32 bytes ok", func(t *testing.T) {
		_, err := NewAESGCMCipher(bytes.Repeat([]byte("k"), 32))
		require.NoError(t, err)
	})
	t.Run("non-32 bytes rejected", func(t *testing.T) {
		_, err := NewAESGCMCipher(bytes.Repeat([]byte("k"), 16))
		require.Error(t, err)
	})
}

func TestAESGCM_RoundTrip(t *testing.T) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	require.NoError(t, err)

	c, err := NewAESGCMCipher(key)
	require.NoError(t, err)

	for _, plaintext := range []string{"", "hello", "a totp secret with non-ascii 中文"} {
		t.Run(plaintext, func(t *testing.T) {
			encoded, err := c.EncryptToString(plaintext)
			require.NoError(t, err)
			require.NotEqual(t, plaintext, encoded)

			decoded, err := c.DecryptFromString(encoded)
			require.NoError(t, err)
			require.Equal(t, plaintext, decoded)
		})
	}
}

func TestAESGCM_NonceUniqueness(t *testing.T) {
	c, err := NewAESGCMCipher(bytes.Repeat([]byte("k"), 32))
	require.NoError(t, err)

	a, err := c.EncryptToString("same-input")
	require.NoError(t, err)
	b, err := c.EncryptToString("same-input")
	require.NoError(t, err)
	require.NotEqual(t, a, b, "ciphertexts must differ due to random nonce")
}

func TestAESGCM_DecryptInvalid(t *testing.T) {
	c, err := NewAESGCMCipher(bytes.Repeat([]byte("k"), 32))
	require.NoError(t, err)

	_, err = c.DecryptFromString("not-base64!!!")
	require.Error(t, err)

	_, err = c.DecryptFromString("AAAA")
	require.Error(t, err)
}
