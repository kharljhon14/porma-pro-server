package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJWTMaker(t *testing.T) {
	secret := "fX7pL2wqE9vB1mZsKj4YtNcRx6HgQeAa"

	jwtMaker, err := NewJWTMaker(secret)
	require.NoError(t, err)
	require.NotEmpty(t, jwtMaker)

	jwtMakerErr, err := NewJWTMaker("12314")
	require.Error(t, err)
	require.Empty(t, jwtMakerErr)
}
