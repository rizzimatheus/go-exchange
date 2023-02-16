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
	transfer1 := createRandomTransfer(t, account1, account2)
	transfer2 := createRandomTransfer(t, account2, account1)

	arg := CreateTradeParams{
		FirstTransferID:  transfer1.ID,
		SecondTransferID: transfer2.ID,
	}

	trade, err := testQueries.CreateTrade(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trade)

	require.Equal(t, arg.FirstTransferID, trade.FirstTransferID)
	require.Equal(t, arg.SecondTransferID, trade.SecondTransferID)

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
	require.Equal(t, trade1.FirstTransferID, trade2.FirstTransferID)
	require.Equal(t, trade1.SecondTransferID, trade2.SecondTransferID)
	require.WithinDuration(t, trade1.CreatedAt, trade2.CreatedAt, time.Second)
}

func TestListTrades(t *testing.T) {
	var lastTrade Trade
	for i := 0; i < 5; i++ {
		lastTrade = createRandomTrade(t)
	}

	arg := ListTradesParams{
		FirstTransferID:  lastTrade.FirstTransferID,
		SecondTransferID: lastTrade.FirstTransferID,
		Limit:            5,
		Offset:           0,
	}

	trades, err := testQueries.ListTrades(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trades)

	for _, trade := range trades {
		require.NotEmpty(t, trade)
		require.Contains(t, []int64{trade.FirstTransferID, trade.SecondTransferID}, arg.FirstTransferID)
	}
}
