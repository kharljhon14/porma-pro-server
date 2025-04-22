package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kharljhon14/porma-pro-server/internal/util"
	"github.com/stretchr/testify/require"
)

func createTestPersonalInfo(t *testing.T) PersonalInfo {
	account := createTestAccount(t)

	args := CreatePersonalInfoParams{
		AccountID:   account.ID,
		Email:       util.RandomEmail(),
		FullName:    util.RandomString(12),
		PhoneNumber: "+639456543438",
		Country:     "Philippines",
		State:       "Bataan",
		City:        "Orion",
	}

	personalInfo, err := testStore.CreatePersonalInfo(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, personalInfo)

	require.Equal(t, args.AccountID, personalInfo.AccountID)
	require.Equal(t, args.Email, personalInfo.Email)
	require.Equal(t, args.FullName, personalInfo.FullName)
	require.Equal(t, args.PhoneNumber, personalInfo.PhoneNumber)
	require.Equal(t, args.Country, personalInfo.Country)
	require.Equal(t, args.State, personalInfo.State)
	require.Equal(t, args.City, personalInfo.City)
	require.Empty(t, personalInfo.LinkedinUrl)
	require.Empty(t, personalInfo.PersonalUrl)

	return personalInfo
}

func TestCreatePersonalInfo(t *testing.T) {
	createTestPersonalInfo(t)
}

func TestGetPersonalInfo(t *testing.T) {
	personalInfo := createTestPersonalInfo(t)

	gotPersonalInfo, err := testStore.GetPersonalInfo(context.Background(), personalInfo.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotPersonalInfo)

	require.Equal(t, personalInfo, gotPersonalInfo)
}

func TestUpdatePersonalInfo(t *testing.T) {
	personalInfo := createTestPersonalInfo(t)

	args := UpdatePersonalInfoParams{
		FullName:    personalInfo.FullName,
		Email:       personalInfo.Email,
		PhoneNumber: personalInfo.PhoneNumber,
		Country:     personalInfo.Country,
		State:       personalInfo.State,
		City:        personalInfo.City,
		ID:          personalInfo.ID,
		LinkedinUrl: pgtype.Text{
			String: "https://www.google.com",
			Valid:  true,
		},
		PersonalUrl: pgtype.Text{
			String: "https://www.google.com",
			Valid:  true,
		},
	}

	updatedPersonalInfo, err := testStore.UpdatePersonalInfo(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, updatedPersonalInfo)

	require.NotEqual(t, personalInfo, updatedPersonalInfo)
	require.NotEmpty(t, updatedPersonalInfo.LinkedinUrl)
	require.NotEmpty(t, updatedPersonalInfo.PersonalUrl)
}

func TestDeletePersonalInfo(t *testing.T) {
	personalInfo := createTestPersonalInfo(t)

	err := testStore.DeletePersonalInfo(context.Background(), personalInfo.ID)
	require.NoError(t, err)

	gotPersonalInfo, err := testStore.GetPersonalInfo(context.Background(), personalInfo.ID)
	require.Error(t, err)
	require.Empty(t, gotPersonalInfo)
}
