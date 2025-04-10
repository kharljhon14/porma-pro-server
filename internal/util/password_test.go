package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "@Password123"

	hashed_password, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashed_password)

	err = CheckPassword(password, hashed_password)
	require.NoError(t, err)

	wrongPassword := "wrongPassword"
	err = CheckPassword(wrongPassword, hashed_password)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
}
