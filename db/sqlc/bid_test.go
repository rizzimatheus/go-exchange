package db

import (
	"context"
	"go-exchange/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomBid(t *testing.T) Bid {
	pair := util.RandomPair()
	c1, c2 := util.CurrenciesFromPair(pair)
	amount := util.RandomMoney()
	account1 := createRandomAccount(t, c1)
	account2 := createRandomAccount(t, c2)

	arg := CreateBidParams{
		Pair:            pair,
		FromAccountID:   account1.ID,
		ToAccountID:     account2.ID,
		Price:           util.RandomMoney(),
		InitialAmount:   amount,
		RemainingAmount: amount,
		Status:          util.RandomStatus(),
	}

	bid, err := testQueries.CreateBid(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, bid)

	require.Equal(t, arg.Pair, bid.Pair)
	require.Equal(t, arg.FromAccountID, bid.FromAccountID)
	require.Equal(t, arg.ToAccountID, bid.ToAccountID)
	require.Equal(t, arg.Price, bid.Price)
	require.Equal(t, arg.InitialAmount, bid.InitialAmount)
	require.Equal(t, arg.RemainingAmount, bid.RemainingAmount)
	require.Equal(t, arg.Status, bid.Status)

	require.NotZero(t, bid.ID)
	require.NotZero(t, bid.CreatedAt)

	return bid
}

func TestCreateBid(t *testing.T) {
	createRandomBid(t)
}

func TestGetBid(t *testing.T) {
	bid1 := createRandomBid(t)
	bid2, err := testQueries.GetBid(context.Background(), bid1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, bid2)

	require.Equal(t, bid1.ID, bid2.ID)
	require.Equal(t, bid1.Pair, bid2.Pair)
	require.Equal(t, bid1.FromAccountID, bid2.FromAccountID)
	require.Equal(t, bid1.ToAccountID, bid2.ToAccountID)
	require.Equal(t, bid1.Price, bid2.Price)
	require.Equal(t, bid1.InitialAmount, bid2.InitialAmount)
	require.Equal(t, bid1.RemainingAmount, bid2.RemainingAmount)
	require.Equal(t, bid1.Status, bid2.Status)
	require.WithinDuration(t, bid1.CreatedAt, bid2.CreatedAt, time.Second)
}

func TestListBids(t *testing.T) {
	var lastBid Bid
	for i := 0; i < 5; i++ {
		lastBid = createRandomBid(t)
	}

	arg := ListBidsParams{
		FromAccountID: lastBid.FromAccountID,
		ToAccountID:   lastBid.FromAccountID,
		Limit:         5,
		Offset:        0,
	}

	bids, err := testQueries.ListBids(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, bids)

	for _, bid := range bids {
		require.NotEmpty(t, bid)
		require.Contains(t, []int64{bid.FromAccountID, bid.ToAccountID}, arg.FromAccountID)
	}
}
