package api

import (
	"bytes"
	"encoding/json"
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
