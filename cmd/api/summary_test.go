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

func TestCreateSummary(t *testing.T) {
	args := createSummmaryRequest{
		AccountID: 1,
		Summary:   util.RandomString(2000),
	}

	summary := db.Summary{
		AccountID: args.AccountID,
		Summary:   args.Summary,
	}

	testCases := []struct {
		name          string
		args          createSummmaryRequest
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			args: args,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					CreateSummary(gomock.Any(), gomock.Eq(db.CreateSummaryParams{
						AccountID: args.AccountID,
						Summary:   args.Summary,
					})).
					Times(1).
					Return(summary, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotSummary db.Summary
				err = json.Unmarshal(data, &gotSummary)
				require.NoError(t, err)

				require.NotEmpty(t, gotSummary)
				require.Equal(t, summary, gotSummary)
			},
		},
		{
			name: "BadRequest",
			args: createSummmaryRequest{
				AccountID: 0,
				Summary:   "",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					CreateSummary(gomock.Any(), gomock.Eq(db.CreateSummaryParams{
						AccountID: 0,
						Summary:   "",
					})).
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
					CreateSummary(gomock.Any(), gomock.Eq(db.CreateSummaryParams{
						AccountID: args.AccountID,
						Summary:   args.Summary,
					})).
					Times(1).
					Return(db.Summary{}, sql.ErrConnDone)
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

			request, err := http.NewRequest(http.MethodPost, "/summary", bytes.NewBuffer(js))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestGetSummary(t *testing.T) {
	id := int64(1)

	summary := db.Summary{
		ID:        id,
		AccountID: 1,
		Summary:   util.RandomString(2000),
	}

	testCases := []struct {
		name          string
		id            int64
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			id:   id,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetSummary(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(summary, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotSummary db.Summary
				err = json.Unmarshal(data, &gotSummary)
				require.NoError(t, err)

				require.NotEmpty(t, gotSummary)
				require.Equal(t, summary, gotSummary)
			},
		},
		{
			name: "BadRequest",
			id:   0,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetSummary(gomock.Any(), gomock.Eq(0)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   id,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetSummary(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(db.Summary{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/summary/%d", tc.id)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateSummary(t *testing.T) {
	args := db.UpdateSummaryParams{
		ID:      1,
		Summary: util.RandomString(1000),
	}

	summary := db.Summary{
		ID:        args.ID,
		AccountID: 1,
		Summary:   args.Summary,
	}

	testCases := []struct {
		name          string
		args          db.UpdateSummaryParams
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			args: args,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					UpdateSummary(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(summary, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotSummary db.Summary
				err = json.Unmarshal(data, &gotSummary)
				require.NoError(t, err)

				require.Equal(t, summary, gotSummary)
			},
		},
		{
			name: "BadRequest",
			args: db.UpdateSummaryParams{},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					UpdateSummary(gomock.Any(), gomock.Eq(db.UpdateSummaryParams{})).
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
					UpdateSummary(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(db.Summary{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/summary/%d", tc.args.ID)

			js, err := json.Marshal(tc.args)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(js))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}

func TestDeleteSummary(t *testing.T) {
	id := int64(1)

	testCases := []struct {
		name          string
		id            int64
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			id:   id,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					DeleteSummary(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			id:   0,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					DeleteSummary(gomock.Any(), gomock.Eq(0)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			id:   id,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					DeleteSummary(gomock.Any(), gomock.Eq(id)).
					Times(1).
					Return(sql.ErrConnDone)
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

			url := fmt.Sprintf("/summary/%d", tc.id)

			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

		})
	}
}
