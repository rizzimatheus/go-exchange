package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "go-exchange/db/mock"
	db "go-exchange/db/sqlc"
	"go-exchange/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateTradeAPI(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)

	account1 := randomAccount(user1.Username)
	account2 := randomAccount(user2.Username)
	account3 := randomAccount(user1.Username)
	account4 := randomAccount(user2.Username)

	account1.Currency = util.BTC
	account2.Currency = util.BTC
	account3.Currency = util.USDT
	account4.Currency = util.USDT

	pair := util.BTC_USDT

	amount := int64(10)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"first_from_account_id":  account1.ID,
				"first_to_account_id":    account2.ID,
				"first_amount":           amount,
				"second_from_account_id": account3.ID,
				"second_to_account_id":   account4.ID,
				"second_amount":          amount,
				"pair":                   pair,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account4.ID)).Times(1).Return(account4, nil)

				arg := db.TradeTxParams{
					FirstFromAccountID: account1.ID,
					FirstToAccountID:   account2.ID,
					FirstAmount:        amount,

					SecondFromAccountID: account3.ID,
					SecondToAccountID:   account4.ID,
					SecondAmount:        amount,
				}
				store.EXPECT().TradeTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "GetAccountError",
			body: gin.H{
				"first_from_account_id":  account1.ID,
				"first_to_account_id":    account2.ID,
				"first_amount":           amount,
				"second_from_account_id": account3.ID,
				"second_to_account_id":   account4.ID,
				"second_amount":          amount,
				"pair":                   pair,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().TradeTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TradeTxError",
			body: gin.H{
				"first_from_account_id":  account1.ID,
				"first_to_account_id":    account2.ID,
				"first_amount":           amount,
				"second_from_account_id": account3.ID,
				"second_to_account_id":   account4.ID,
				"second_amount":          amount,
				"pair":                   pair,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account4.ID)).Times(1).Return(account4, nil)

				store.EXPECT().TradeTx(gomock.Any(), gomock.Any()).Times(1).Return(db.TradeTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPair",
			body: gin.H{
				"first_from_account_id":  account1.ID,
				"first_to_account_id":    account2.ID,
				"first_amount":           amount,
				"second_from_account_id": account3.ID,
				"second_to_account_id":   account4.ID,
				"second_amount":          amount,
				"pair":                   "BTC/BTC",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(0).Return(account1, nil)
				store.EXPECT().TradeTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server, err := NewServer(util.Config{}, store)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/trades"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
