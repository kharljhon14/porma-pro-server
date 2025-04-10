package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/kharljhon14/porma-pro-server/internal/util"
	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T) Account {
	hashedPassword, err := util.HashedPassword(util.RandomString(8))
	require.NoError(t, err)

	args := CreateAccountParams{
		Email:        util.RandomEmail(),
		PasswordHash: hashedPassword,
		FullName:     util.RandomString(10),
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

func TestGetAccount(t *testing.T) {
	account := createTestAccount(t)

	account2, err := testStore.GetAccount(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account.FullName, account2.FullName)
	require.Equal(t, account.Email, account2.Email)
	require.Equal(t, account.PasswordHash, account2.PasswordHash)
	require.WithinDuration(t, account.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createTestAccount(t)

	err := testStore.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account.ID)
	require.Error(t, err)
	errorMessage := fmt.Errorf("sql: %s", err)
	require.EqualError(t, errorMessage, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}
