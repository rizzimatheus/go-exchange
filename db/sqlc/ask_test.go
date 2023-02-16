package db

import (
	"context"
	"go-exchange/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAsk(t *testing.T) Ask {
	pair := util.RandomPair()
	c1, c2 := util.CurrenciesFromPair(pair)
	account1 := createRandomAccount(t, c1)
	account2 := createRandomAccount(t, c2)

	arg := CreateAskParams{
		Pair:          pair,
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Price:         util.RandomMoney(),
		Amount:        util.RandomMoney(),
		Status:        util.RandomStatus(),
	}

	ask, err := testQueries.CreateAsk(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, ask)

	require.Equal(t, arg.Pair, ask.Pair)
	require.Equal(t, arg.FromAccountID, ask.FromAccountID)
	require.Equal(t, arg.ToAccountID, ask.ToAccountID)
	require.Equal(t, arg.Price, ask.Price)
	require.Equal(t, arg.Amount, ask.Amount)
	require.Equal(t, arg.Status, ask.Status)

	require.NotZero(t, ask.ID)
	require.NotZero(t, ask.CreatedAt)

	return ask
}

func TestCreateAsk(t *testing.T) {
	createRandomAsk(t)
}

func TestGetAsk(t *testing.T) {
	ask1 := createRandomAsk(t)
	ask2, err := testQueries.GetAsk(context.Background(), ask1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, ask2)

	require.Equal(t, ask1.ID, ask2.ID)
	require.Equal(t, ask1.Pair, ask2.Pair)
	require.Equal(t, ask1.FromAccountID, ask2.FromAccountID)
	require.Equal(t, ask1.ToAccountID, ask2.ToAccountID)
	require.Equal(t, ask1.Price, ask2.Price)
	require.Equal(t, ask1.Amount, ask2.Amount)
	require.Equal(t, ask1.Status, ask2.Status)
	require.WithinDuration(t, ask1.CreatedAt, ask2.CreatedAt, time.Second)
}

func TestListAsks(t *testing.T) {
	var lastAsk Ask
	for i := 0; i < 5; i++ {
		lastAsk = createRandomAsk(t)
	}

	arg := ListAsksParams{
		FromAccountID: lastAsk.FromAccountID,
		ToAccountID:   lastAsk.FromAccountID,
		Limit:         5,
		Offset:        0,
	}

	asks, err := testQueries.ListAsks(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, asks)

	for _, ask := range asks {
		require.NotEmpty(t, ask)
		require.Contains(t, []int64{arg.FromAccountID, arg.ToAccountID}, ask.FromAccountID)
	}
}
