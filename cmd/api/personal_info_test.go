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

	"github.com/golang/mock/gomock"
	mock_sqlc "github.com/kharljhon14/porma-pro-server/internal/db/mock"
	db "github.com/kharljhon14/porma-pro-server/internal/db/sqlc"
	"github.com/kharljhon14/porma-pro-server/internal/util"
	"github.com/stretchr/testify/require"
)

func TestCreatePersonalInfo(t *testing.T) {

	args := createPersonalInfoRequest{
		AccountID:   1,
		FullName:    util.RandomString(12),
		Email:       util.RandomEmail(),
		PhoneNumber: "+639456543438",
		Country:     "Philippines",
		State:       "Bataan",
		City:        "Orion",
	}

	peronsalInfo := db.PersonalInfo{
		ID:          util.RandomInt(1, 1000),
		AccountID:   args.AccountID,
		FullName:    args.FullName,
		Email:       args.Email,
		PhoneNumber: args.PhoneNumber,
		Country:     args.Country,
		State:       args.State,
		City:        args.City,
	}

	testCases := []struct {
		name          string
		args          createPersonalInfoRequest
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			args: args,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					CreatePersonalInfo(gomock.Any(), gomock.Any()).
					Times(1).
					Return(peronsalInfo, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotPersonalInfo db.PersonalInfo
				err = json.Unmarshal(data, &gotPersonalInfo)
				require.NoError(t, err)

				require.Equal(t, gotPersonalInfo, peronsalInfo)
			},
		},
		{
			name: "BadRequest",
			args: createPersonalInfoRequest{
				AccountID: 0,
				Email:     "invalid",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					CreatePersonalInfo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			args: args,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					CreatePersonalInfo(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.PersonalInfo{}, sql.ErrConnDone)
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

			store := mock_sqlc.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestingServer(t, store)
			recorder := httptest.NewRecorder()

			js, err := json.Marshal(tc.args)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/personal-info", bytes.NewBuffer(js))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetPersonalInfo(t *testing.T) {
	personalInfo := db.PersonalInfo{
		ID:          1,
		AccountID:   1,
		FullName:    util.RandomString(12),
		Email:       util.RandomEmail(),
		PhoneNumber: "+639456543438",
		Country:     "Philippines",
		State:       "Bataan",
		City:        "Orion",
	}

	testCases := []struct {
		name          string
		args          int64
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			args: personalInfo.ID,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetPersonalInfo(gomock.Any(), gomock.Eq(personalInfo.ID)).
					Times(1).
					Return(personalInfo, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotPersonalInfo db.PersonalInfo
				err = json.Unmarshal(data, &gotPersonalInfo)
				require.NoError(t, err)

				require.Equal(t, personalInfo, gotPersonalInfo)
			},
		},
		{
			name: "BadRequest",
			args: 0,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetPersonalInfo(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			args: personalInfo.ID,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetPersonalInfo(gomock.Any(), gomock.Eq(personalInfo.ID)).
					Times(1).
					Return(db.PersonalInfo{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			args: personalInfo.ID,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetPersonalInfo(gomock.Any(), gomock.Eq(personalInfo.ID)).
					Times(1).
					Return(db.PersonalInfo{}, sql.ErrConnDone)
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

			store := mock_sqlc.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestingServer(t, store)

			url := fmt.Sprintf(`/personal-info/%d`, tc.args)

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
