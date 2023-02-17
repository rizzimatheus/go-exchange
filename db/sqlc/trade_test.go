package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomTrade(t *testing.T) Trade {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	account3 := createRandomAccount(t)
	account4 := createRandomAccount(t)

	transfer1 := createRandomTransfer(t, account1, account2)
	transfer2 := createRandomTransfer(t, account3, account4)

	arg := CreateTradeParams{
		FirstFromAccountID:  transfer1.FromAccountID,
		FirstToAccountID:    transfer1.ToAccountID,
		FirstAmount:         transfer1.Amount,
		SecondFromAccountID: transfer2.FromAccountID,
		SecondToAccountID:   transfer2.ToAccountID,
		SecondAmount:        transfer2.Amount,
	}

	trade, err := testQueries.CreateTrade(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trade)

	require.Equal(t, arg.FirstFromAccountID, trade.FirstFromAccountID)
	require.Equal(t, arg.FirstToAccountID, trade.FirstToAccountID)
	require.Equal(t, arg.FirstAmount, trade.FirstAmount)

	require.Equal(t, arg.SecondFromAccountID, trade.SecondFromAccountID)
	require.Equal(t, arg.SecondToAccountID, trade.SecondToAccountID)
	require.Equal(t, arg.SecondAmount, trade.SecondAmount)

	require.NotZero(t, trade.ID)
	require.NotZero(t, trade.CreatedAt)

	return trade
}

func TestCreateTrade(t *testing.T) {
	createRandomTrade(t)
}

func TestGetTrade(t *testing.T) {
	trade1 := createRandomTrade(t)
	trade2, err := testQueries.GetTrade(context.Background(), trade1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, trade2)

	require.Equal(t, trade1.ID, trade2.ID)

	require.Equal(t, trade1.FirstFromAccountID, trade2.FirstFromAccountID)
	require.Equal(t, trade1.FirstToAccountID, trade2.FirstToAccountID)
	require.Equal(t, trade1.FirstAmount, trade2.FirstAmount)

	require.Equal(t, trade1.SecondFromAccountID, trade2.SecondFromAccountID)
	require.Equal(t, trade1.SecondToAccountID, trade2.SecondToAccountID)
	require.Equal(t, trade1.SecondAmount, trade2.SecondAmount)

	require.WithinDuration(t, trade1.CreatedAt, trade2.CreatedAt, time.Second)
}

func TestListTrades(t *testing.T) {
	var lastTrade Trade
	for i := 0; i < 5; i++ {
		lastTrade = createRandomTrade(t)
	}

	arg := ListTradesParams{
		FirstFromAccountID: lastTrade.FirstFromAccountID,
		FirstToAccountID: lastTrade.FirstFromAccountID,
		SecondFromAccountID: lastTrade.FirstFromAccountID,
		SecondToAccountID: lastTrade.FirstFromAccountID,
		Limit:            5,
		Offset:           0,
	}

	trades, err := testQueries.ListTrades(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trades)

	for _, trade := range trades {
		require.NotEmpty(t, trade)
		require.Contains(
			t, 
			[]int64{trade.FirstFromAccountID, trade.FirstToAccountID, trade.SecondFromAccountID, trade.SecondToAccountID}, 
			arg.FirstFromAccountID)
	}
}
