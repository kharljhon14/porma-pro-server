package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5/pgtype"
	mock_db "github.com/kharljhon14/porma-pro-server/internal/db/mock"
	db "github.com/kharljhon14/porma-pro-server/internal/db/sqlc"
	"github.com/kharljhon14/porma-pro-server/internal/util"
	"github.com/stretchr/testify/require"
)

func TestCreateWorkExperience(t *testing.T) {

	args := createWorkExperienceRequest{
		AccountID: 1,
		Role:      "Web Developer",
		Company:   "KharlDEV",
		Location:  "Philippines",
		Summary:   util.RandomString(10),
		StartDate: time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
	}

	workExperience := db.WorkExperience{
		ID:        1,
		AccountID: args.AccountID,
		Role:      args.Role,
		Company:   args.Company,
		Location:  args.Location,
		Summary:   args.Summary,
		StartDate: pgtype.Timestamp{
			Valid: true,
			Time:  args.StartDate,
		},
		EndDate: pgtype.Timestamp{
			Valid: true,
			Time:  args.EndDate,
		},
	}

	testCases := []struct {
		name          string
		args          createWorkExperienceRequest
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			args: args,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					CreateWorkExperience(gomock.Any(), gomock.Eq(db.CreateWorkExperienceParams{
						AccountID: args.AccountID,
						Role:      args.Role,
						Company:   args.Company,
						Location:  args.Location,
						Summary:   args.Summary,
						StartDate: pgtype.Timestamp{
							Valid: true,
							Time:  args.StartDate,
						},
						EndDate: pgtype.Timestamp{
							Valid: true,
							Time:  args.EndDate,
						},
					})).
					Times(1).
					Return(workExperience, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotWorkExperience db.WorkExperience
				err = json.Unmarshal(data, &gotWorkExperience)
				require.NoError(t, err)

				require.NotEmpty(t, gotWorkExperience)
				require.Equal(t, workExperience, gotWorkExperience)
			},
		},
		{
			name: "BadRequest",
			args: createWorkExperienceRequest{},
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					CreateWorkExperience(gomock.Any(), gomock.Eq(db.CreateWorkExperienceParams{})).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			args: args,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					CreateWorkExperience(gomock.Any(), gomock.Eq(db.CreateWorkExperienceParams{
						AccountID: args.AccountID,
						Role:      args.Role,
						Company:   args.Company,
						Location:  args.Location,
						Summary:   args.Summary,
						StartDate: pgtype.Timestamp{
							Valid: true,
							Time:  args.StartDate,
						},
						EndDate: pgtype.Timestamp{
							Valid: true,
							Time:  args.EndDate,
						},
					})).
					Times(1).
					Return(db.WorkExperience{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestingServer(t, store)

			recorder := httptest.NewRecorder()

			js, err := json.Marshal(tc.args)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/work-experience", bytes.NewBuffer(js))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetWorkExperience(t *testing.T) {

	id := int64(1)

	workExperience := db.WorkExperience{
		ID:        id,
		AccountID: 1,
		Role:      "Dev",
		Company:   "Test",
		Location:  "Bataan",
		Summary:   util.RandomString(12),
		StartDate: pgtype.Timestamp{
			Valid: true,
			Time:  time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC),
		},
		EndDate: pgtype.Timestamp{
			Valid: true,
			Time:  time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	testCases := []struct {
		name          string
		args          int64
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			args: id,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					GetWorkExperience(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(workExperience, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotWorkExperience db.WorkExperience
				err = json.Unmarshal(data, &gotWorkExperience)
				require.NoError(t, err)

				require.NotEmpty(t, gotWorkExperience)
				require.Equal(t, workExperience, gotWorkExperience)

			},
		},
		{
			name: "BadRequest",
			args: 0,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					GetWorkExperience(gomock.Any(), gomock.Eq(0)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			args: id,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					GetWorkExperience(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(db.WorkExperience{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			args: id,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					GetWorkExperience(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(db.WorkExperience{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestingServer(t, store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/work-experience/%d", tc.args)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}

}

func TestGetWorkExperienceList(t *testing.T) {

	id := int64(1)

	workExperiences := []db.WorkExperience{
		{
			ID:        1,
			AccountID: id,
			Role:      "Dev",
			Company:   "Test",
			Location:  "Bataan",
			Summary:   util.RandomString(12),
			StartDate: pgtype.Timestamp{
				Valid: true,
				Time:  time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC),
			},
			EndDate: pgtype.Timestamp{
				Valid: true,
				Time:  time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			ID:        2,
			AccountID: id,
			Role:      "Dev2",
			Company:   "Test2",
			Location:  "Bataan2",
			Summary:   util.RandomString(12),
			StartDate: pgtype.Timestamp{
				Valid: true,
				Time:  time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC),
			},
			EndDate: pgtype.Timestamp{
				Valid: true,
				Time:  time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	testCases := []struct {
		name          string
		args          int64
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			args: id,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					GetWorkExperiences(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(workExperiences, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotWorkExperiences []db.WorkExperience
				err = json.Unmarshal(data, &gotWorkExperiences)
				require.NoError(t, err)

				require.Equal(t, len(workExperiences), len(gotWorkExperiences))
				require.Equal(t, workExperiences[0], gotWorkExperiences[0])
				require.Equal(t, workExperiences[1], gotWorkExperiences[1])
			},
		},
		{
			name: "BadRequest",
			args: 0,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					GetWorkExperiences(gomock.Any(), gomock.Eq(0)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			args: id,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					GetWorkExperiences(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return([]db.WorkExperience{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			args: id,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					GetWorkExperiences(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return([]db.WorkExperience{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestingServer(t, store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/work-experience/?account_id=%d", tc.args)

			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateWorkExperience(t *testing.T) {
	id := int64(1)

	args := updateWorkExperienceRequest{
		Role:      "Web Developer",
		Company:   "KharlDEV",
		Location:  "Philippines",
		Summary:   util.RandomString(10),
		StartDate: time.Date(2024, time.April, 2, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, time.April, 2, 0, 0, 0, 0, time.UTC),
	}

	workExperience := db.WorkExperience{
		ID:        id,
		AccountID: 1,
		Role:      args.Role,
		Company:   args.Company,
		Location:  args.Location,
		Summary:   args.Summary,
		StartDate: pgtype.Timestamp{
			Valid: true,
			Time:  args.StartDate,
		},
		EndDate: pgtype.Timestamp{
			Valid: true,
			Time:  args.EndDate,
		},
	}

	testCases := []struct {
		name          string
		id            int64
		args          updateWorkExperienceRequest
		buildStubs    func(store *mock_db.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			id:   id,
			args: args,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					UpdateWorkExperience(gomock.Any(), gomock.Eq(db.UpdateWorkExperienceParams{
						ID:       id,
						Role:     args.Role,
						Company:  args.Company,
						Location: args.Location,
						Summary:  args.Summary,
						StartDate: pgtype.Timestamp{
							Valid: true,
							Time:  args.StartDate,
						},
						EndDate: pgtype.Timestamp{
							Valid: true,
							Time:  args.EndDate,
						},
					})).
					Times(1).
					Return(workExperience, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotWorkExperience db.WorkExperience
				err = json.Unmarshal(data, &gotWorkExperience)
				require.NoError(t, err)

				require.NotEmpty(t, gotWorkExperience)
				require.Equal(t, workExperience, gotWorkExperience)
			},
		},
		{
			name: "BadRequest",
			id:   0,
			args: updateWorkExperienceRequest{},
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					UpdateWorkExperience(gomock.Any(), gomock.Eq(db.UpdateWorkExperienceParams{})).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   id,
			args: args,
			buildStubs: func(store *mock_db.MockStore) {
				store.
					EXPECT().
					UpdateWorkExperience(gomock.Any(), gomock.Eq(db.UpdateWorkExperienceParams{
						ID:       id,
						Role:     args.Role,
						Company:  args.Company,
						Location: args.Location,
						Summary:  args.Summary,
						StartDate: pgtype.Timestamp{
							Valid: true,
							Time:  args.StartDate,
						},
						EndDate: pgtype.Timestamp{
							Valid: true,
							Time:  args.EndDate,
						},
					})).
					Times(1).
					Return(db.WorkExperience{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mock_db.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestingServer(t, store)

			recorder := httptest.NewRecorder()

			js, err := json.Marshal(tc.args)
			require.NoError(t, err)

			url := fmt.Sprintf("/work-experience/%d", tc.id)

			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(js))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}
