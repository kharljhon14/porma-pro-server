package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
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
		ID:           util.RandomgInt(1, 1000),
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
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
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
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recoder.Code)

				data, err := io.ReadAll(recoder.Body)
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
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recoder.Code)
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
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recoder.Code)
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
