package db

import (
	"context"
	"testing"

	"github.com/kharljhon14/porma-pro-server/internal/util"
	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T) Account {
	hashedPassword, err := util.HashedPassword("Password")
	require.NoError(t, err)

	args := CreateAccountParams{
		Email:        "enriquezkharl14@gmail.com",
		PasswordHash: hashedPassword,
		FullName:     "Kharl Jhon Rhane Enriquez",
	}

	account, err := testStore.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.Email, account.Email)
	require.Equal(t, args.PasswordHash, account.PasswordHash)
	require.Equal(t, args.FullName, account.FullName)
	require.Equal(t, false, account.IsVerified)

	require.NotZero(t, account.CreatedAt)
	require.NotZero(t, account.UpdatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}
