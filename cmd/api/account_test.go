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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	mock_sqlc "github.com/kharljhon14/porma-pro-server/internal/db/mock"
	db "github.com/kharljhon14/porma-pro-server/internal/db/sqlc"
	"github.com/kharljhon14/porma-pro-server/internal/util"
	"github.com/stretchr/testify/require"
)

func TestCreateAccountAPI(t *testing.T) {

	args := createAccountRequest{
		Email:    util.RandomEmail(),
		Password: util.RandomString(12),
		FullName: util.RandomString(12),
	}

	account := db.Account{
		ID:           util.RandomInt(1, 1000),
		Email:        args.Email,
		PasswordHash: args.Password,
		FullName:     args.FullName,
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		IsVerified: false,
	}

	testCases := []struct {
		name          string
		args          createAccountRequest
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Created",
			args: args,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotAccount db.Account
				err = json.Unmarshal(data, &gotAccount)
				require.NoError(t, err)

				require.Equal(t, gotAccount, gotAccount)
			},
		},
		{
			name: "BadRequest",
			args: createAccountRequest{},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(createAccountRequest{})).
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
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Forbidden",
			args: createAccountRequest{
				Email:    "enriquezkharl14@gmail.com",
				Password: "@paswword123",
				FullName: "Kharl Curz",
			},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, &pgconn.PgError{
						Code:           "23505",
						ConstraintName: "accounts_email_key",
					})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
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

			url := "/sign-up"

			js, err := json.Marshal(tc.args)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(js))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestLoginAccountAPI(t *testing.T) {
	args := loginAccountRequest{
		Email:    util.RandomEmail(),
		Password: "@Password123",
	}

	hashedPassword, err := util.HashedPassword("@Password123")
	require.NoError(t, err)

	account := db.Account{
		ID:           util.RandomInt(1, 1000),
		Email:        args.Email,
		FullName:     util.RandomString(12),
		PasswordHash: hashedPassword,
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		IsVerified: true,
	}

	testCases := []struct {
		name          string
		args          loginAccountRequest
		buildStubs    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder httptest.ResponseRecorder)
	}{
		{
			name: "Ok",
			args: args,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetAccountByEmail(gomock.Any(), gomock.Eq(args.Email)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			args: loginAccountRequest{},
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetAccountByEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NotFound",
			args: args,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetAccountByEmail(gomock.Any(), gomock.Eq(args.Email)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			args: args,
			buildStubs: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetAccountByEmail(gomock.Any(), gomock.Eq(args.Email)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
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

			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(js))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, *recorder)
		})
	}
}

func TestGetAccountAPI(t *testing.T) {
	hashed_password, err := util.HashedPassword("@Password123")
	require.NoError(t, err)

	account := db.Account{
		ID:           util.RandomInt(1, 1000),
		Email:        util.RandomEmail(),
		PasswordHash: hashed_password,
		FullName:     util.RandomString(12),
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		IsVerified: false,
	}

	testCases := []struct {
		name          string
		accountID     int64
		buildStuds    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "Ok",
			accountID: account.ID,
			buildStuds: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotAccount db.Account
				err = json.Unmarshal(data, &gotAccount)
				require.NoError(t, err)

				require.Equal(t, account, gotAccount)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStuds: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			buildStuds: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStuds: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
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
			tc.buildStuds(store)

			server := newTestingServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestVerifyAccountAPI(t *testing.T) {
	hashed_password, err := util.HashedPassword("@Password123")
	require.NoError(t, err)

	account := db.Account{
		ID:           util.RandomInt(1, 1000),
		Email:        util.RandomEmail(),
		PasswordHash: hashed_password,
		FullName:     util.RandomString(12),
		CreatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		IsVerified: false,
	}

	testCases := []struct {
		name          string
		accountID     int64
		buildStuds    func(store *mock_sqlc.MockStore)
		checkResponse func(t *testing.T, recorder httptest.ResponseRecorder)
	}{
		{
			name:      "Ok",
			accountID: account.ID,
			buildStuds: func(store *mock_sqlc.MockStore) {
				store.EXPECT().
					VerifyAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{
						ID:           account.ID,
						Email:        account.Email,
						PasswordHash: hashed_password,
						FullName:     account.FullName,
						CreatedAt:    account.CreatedAt,
						UpdatedAt:    account.UpdatedAt,
						IsVerified:   true,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)

				data, err := io.ReadAll(recorder.Body)
				require.NoError(t, err)

				var gotAccount db.Account
				err = json.Unmarshal(data, &gotAccount)
				require.NoError(t, err)

				require.Equal(t, true, gotAccount.IsVerified)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			buildStuds: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					VerifyAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStuds: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					VerifyAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStuds: func(store *mock_sqlc.MockStore) {
				store.
					EXPECT().
					VerifyAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder httptest.ResponseRecorder) {
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
			tc.buildStuds(store)

			server := newTestingServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/verify/%d", tc.accountID)
			req, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)

			tc.checkResponse(t, *recorder)
		})
	}
}
