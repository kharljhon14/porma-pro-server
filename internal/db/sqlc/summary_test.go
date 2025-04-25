package db

import (
	"context"
	"testing"

	"github.com/kharljhon14/porma-pro-server/internal/util"
	"github.com/stretchr/testify/require"
)

func createTestSummary(t *testing.T) Summary {
	account := createTestAccount(t)

	args := CreateSummaryParams{
		AccountID: account.ID,
		Summary:   util.RandomString(1000),
	}

	summary, err := testStore.CreateSummary(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, summary)

	require.Equal(t, args.AccountID, summary.AccountID)
	require.Equal(t, args.Summary, summary.Summary)

	return summary
}

func TestCreateSummary(t *testing.T) {
	createTestSummary(t)
}

func TestGetSummary(t *testing.T) {
	summary := createTestSummary(t)

	gotSummary, err := testStore.GetSummary(context.Background(), summary.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotSummary)
	require.Equal(t, summary, gotSummary)
}

func TestUpdateSummary(t *testing.T) {
	summary := createTestSummary(t)

	args := UpdateSummaryParams{
		Summary: util.RandomString(1002),
		ID:      summary.ID,
	}

	updatedSummary, err := testStore.UpdateSummary(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedSummary)

	require.NotEqual(t, summary, updateSummary)
}

func TestDeleteSummary(t *testing.T) {
	summary := createTestSummary(t)

	err := testStore.DeleteSummary(context.Background(), summary.ID)
	require.NoError(t, err)

	gotSummary, err := testStore.GetSummary(context.Background(), summary.ID)
	require.Error(t, err)
	require.Empty(t, gotSummary)
}
