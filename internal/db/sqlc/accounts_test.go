package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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

func TestGetAccountByEmail(t *testing.T) {
	account := createTestAccount(t)

	account2, err := testStore.GetAccountByEmail(context.Background(), account.Email)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account.FullName, account2.FullName)
	require.Equal(t, account.Email, account2.Email)
	require.Equal(t, account.PasswordHash, account2.PasswordHash)
	require.WithinDuration(t, account.CreatedAt.Time, account2.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createTestAccount(t)

	args := UpdateAccountParams{
		FullName: util.RandomString(12),
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		ID: account.ID,
	}
	account2, err := testStore.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, args.FullName, account2.FullName)
	require.WithinDuration(t, args.UpdatedAt.Time, account2.UpdatedAt.Time, time.Second)
}

func TestVerifyAccount(t *testing.T) {
	account := createTestAccount(t)

	account2, err := testStore.VerifyAccount(context.Background(), account.ID)
	require.NoError(t, err)

	require.Equal(t, true, account2.IsVerified)
}

func TestDeleteAccount(t *testing.T) {
	account := createTestAccount(t)

	err := testStore.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account.ID)
	require.Error(t, err)

	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, account2)
}
