package db

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kharljhon14/porma-pro-server/internal/util"
	"github.com/stretchr/testify/require"
)

func createTestWorkExperience(t *testing.T, account Account) WorkExperience {

	args := CreateWorkExperienceParams{
		AccountID: account.ID,
		Role:      "Developer",
		Company:   "KarlDEV",
		Location:  "Philippines",
		Summary:   util.RandomString(1000),
		StartDate: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
		EndDate: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	}

	workExp, err := testStore.CreateWorkExperience(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, workExp)

	return workExp
}

func TestCreateWorkExperience(t *testing.T) {
	account := createTestAccount(t)
	createTestWorkExperience(t, account)
}

func TestGetWorkExperience(t *testing.T) {
	account := createTestAccount(t)
	workExperience := createTestWorkExperience(t, account)

	gotWorkExperience, err := testStore.GetWorkExperience(context.Background(), workExperience.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotWorkExperience)
	require.Equal(t, workExperience, gotWorkExperience)
}

func TestGetWorkExperiences(t *testing.T) {
	account := createTestAccount(t)

	workExperience1 := createTestWorkExperience(t, account)
	workExperience2 := createTestWorkExperience(t, account)

	gotWorkExperiences, err := testStore.GetWorkExperiences(context.Background(), account.ID)
	require.NoError(t, err)
	require.Equal(t, 2, len(gotWorkExperiences))

	require.NotEmpty(t, gotWorkExperiences[0])
	require.NotEmpty(t, gotWorkExperiences[1])

	require.Equal(t, gotWorkExperiences[0], workExperience1)
	require.Equal(t, gotWorkExperiences[1], workExperience2)
}

func TestUpdateWorkExperience(t *testing.T) {
	account := createTestAccount(t)

	workExperience := createTestWorkExperience(t, account)

	args := UpdateWorkExperienceParams{
		Role:      "CEO",
		ID:        workExperience.ID,
		StartDate: workExperience.StartDate,
		EndDate:   workExperience.EndDate,
	}

	updatedWorkExperience, err := testStore.UpdateWorkExperience(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedWorkExperience)

	require.NotEqual(t, workExperience, updatedWorkExperience)
	require.Equal(t, args.Role, updatedWorkExperience.Role)
}

func TestDeleteWorkExperience(t *testing.T) {
	account := createTestAccount(t)

	workExperience := createTestWorkExperience(t, account)

	err := testStore.DeleteWorkExperience(context.Background(), workExperience.ID)
	require.NoError(t, err)

	gotWorkExperience, err := testStore.GetWorkExperience(context.Background(), workExperience.ID)
	require.Error(t, err)
	require.Empty(t, gotWorkExperience)
}
