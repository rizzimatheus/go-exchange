package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "go-exchange/db/mock"
	db "go-exchange/db/sqlc"
	"go-exchange/token"
	"go-exchange/util"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomAsk(fromAccountID int64, toAccountID int64) db.Ask {
	amount := util.RandomInt(1, 99)
	return db.Ask{
		Pair:            util.RandomPair(),
		FromAccountID:   fromAccountID,
		ToAccountID:     toAccountID,
		Price:           util.RandomMoney(),
		InitialAmount:   amount,
		RemainingAmount: amount,
	}
}

func requireBodyMatchAsk(t *testing.T, body *bytes.Buffer, ask db.Ask) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAsk db.Ask
	err = json.Unmarshal(data, &gotAsk)
	require.NoError(t, err)
	require.Equal(t, ask, gotAsk)
}

func TestCreateAskAPI(t *testing.T) {
	user, _ := randomUser(t)

	account1 := randomAccount(user.Username)
	account2 := randomAccount(user.Username)
	account3 := randomAccount(user.Username)

	ask := randomAsk(account1.ID, account2.ID)
	ask.Pair = util.BTC_USDT
	account1.Currency = util.BTC
	account2.Currency = util.USDT
	account3.Currency = util.ETH

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				group := int32(1)
				listArg := db.ListTradableBidsParams{
					Pair:   ask.Pair,
					Price:  ask.Price,
					Limit:  groupSize,
					Offset: (group - 1) * groupSize,
				}
				store.EXPECT().ListTradableBids(gomock.Any(), gomock.Eq(listArg)).Times(1).Return([]db.Bid{}, nil)

				arg := db.CreateAskParams{
					Pair:            ask.Pair,
					FromAccountID:   ask.FromAccountID,
					ToAccountID:     ask.ToAccountID,
					Price:           ask.Price,
					InitialAmount:   ask.InitialAmount,
					RemainingAmount: ask.RemainingAmount,
					Status:          util.ACTIVE,
				}
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Eq(arg)).Times(1).Return(ask, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAsk(t, recorder.Body, ask)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(account2, nil)

				group := int32(1)
				listArg := db.ListTradableBidsParams{
					Pair:   ask.Pair,
					Price:  ask.Price,
					Limit:  groupSize,
					Offset: (group - 1) * groupSize,
				}
				store.EXPECT().ListTradableBids(gomock.Any(), gomock.Eq(listArg)).Times(1).Return([]db.Bid{}, nil)

				arg := db.CreateAskParams{
					Pair:            ask.Pair,
					FromAccountID:   ask.FromAccountID,
					ToAccountID:     ask.ToAccountID,
					Price:           ask.Price,
					InitialAmount:   ask.InitialAmount,
					RemainingAmount: ask.RemainingAmount,
					Status:          util.ACTIVE,
				}

				store.EXPECT().CreateAsk(gomock.Any(), gomock.Eq(arg)).Times(1).Return(db.Ask{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "FromAccountNotFound",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "ToAccountNotFound",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "FromAccountCurrencyMismatch",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account3.ID,
				"to_account_id":   account1.ID,
				"price":           ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "ToAccountCurrencyMismatch",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account3.ID,
				"price":           ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account3.ID)).Times(1).Return(account3, nil)
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPair",
			body: gin.H{
				"pair":            "BTC_BTC",
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativePrice",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           -ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           ask.Price,
				"amount":          -ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccountError",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "UnauthorizedUser",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unauthorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"pair":            ask.Pair,
				"from_account_id": account1.ID,
				"to_account_id":   account2.ID,
				"price":           ask.Price,
				"amount":          ask.InitialAmount,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				//
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account1.ID)).Times(0)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(account2.ID)).Times(0)
				store.EXPECT().CreateAsk(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
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

			url := "/asks"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
