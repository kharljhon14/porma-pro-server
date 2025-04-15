package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPayload(t *testing.T) {
	email := "test@mail.com"
	duration := 24 * time.Hour

	payload, err := NewPayload(email, duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.Equal(t, payload.Email, email)
	require.WithinDuration(
		t,
		payload.IssuedAt,
		time.Now(),
		time.Second,
	)
	require.WithinDuration(
		t,
		payload.ExpiredAt,
		payload.IssuedAt.Add(duration),
		time.Second,
	)
}
