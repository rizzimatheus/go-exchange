package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "go-exchange/db/mock"
	db "go-exchange/db/sqlc"
	"go-exchange/util"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomBid(fromAccountID int64, toAccountID int64) db.Bid {
	return db.Bid{
		Pair:          util.RandomPair(),
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Price:         util.RandomMoney(),
		Amount:        util.RandomInt(1, 99),
	}
}

func requireBodyMatchBid(t *testing.T, body *bytes.Buffer, bid db.Bid) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotBid db.Bid
	err = json.Unmarshal(data, &gotBid)
	require.NoError(t, err)
	require.Equal(t, bid, gotBid)
}

func TestCreateBidAPI(t *testing.T) {
	user, _ := randomUser(t)

	account1 := randomAccount(user.Username)
	account2 := randomAccount(user.Username)
	account3 := randomAccount(user.Username)

	bid := randomBid(account1.ID, account2.ID)
	bid.Pair = util.BTC_USDT
	account1.Currency = util.USDT
	account2.Currency = util.BTC
	account3.Currency = util.ETH

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"pair":            bid.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           bid.Price,
				"amount":          bid.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				arg := db.CreateBidParams{
					Pair:          bid.Pair,
					FromAccountID: bid.FromAccountID,
					ToAccountID:   bid.ToAccountID,
					Price:         bid.Price,
					Amount:        bid.Amount,
					Status:        util.ACTIVE,
				}

				store.EXPECT().CreateBid(gomock.Any(), gomock.Eq(arg)).Times(1).Return(bid, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchBid(t, recorder.Body, bid)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"pair":            bid.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           bid.Price,
				"amount":          bid.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				arg := db.CreateBidParams{
					Pair:          bid.Pair,
					FromAccountID: bid.FromAccountID,
					ToAccountID:   bid.ToAccountID,
					Price:         bid.Price,
					Amount:        bid.Amount,
					Status:        util.ACTIVE,
				}

				store.EXPECT().CreateBid(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Bid{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"pair":            bid.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           bid.Price,
				"amount":          bid.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().CreateBid(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"pair":            bid.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           bid.Price,
				"amount":          bid.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().CreateBid(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			body: gin.H{
				"pair":            bid.Pair,
				"from_account_id": account3.ID,
				"to_account_id":   account1.ID,
				"price":           bid.Price,
				"amount":          bid.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().CreateBid(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: gin.H{
				"pair":            bid.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account3.ID,
				"price":           bid.Price,
				"amount":          bid.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().CreateBid(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPair",
			body: gin.H{
				"pair":            "BTC/BTC",
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           bid.Price,
				"amount":          bid.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateBid(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativePrice",
			body: gin.H{
				"pair":            bid.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           -bid.Price,
				"amount":          bid.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateBid(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"pair":            bid.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           bid.Price,
				"amount":          -bid.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateBid(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccountError",
			body: gin.H{
				"pair":            bid.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           bid.Price,
				"amount":          bid.Amount,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().CreateBid(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/bids"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
