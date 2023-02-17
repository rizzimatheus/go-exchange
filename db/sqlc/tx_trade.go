package db

import "context"

// TradeTxParams contains the input parameters of the trade transaction
type TradeTxParams struct {
	FirstFromAccountID  int64 `json:"first_from_account_id"`
	FirstToAccountID    int64 `json:"first_to_account_id"`
	FirstAmount         int64 `json:"first_amount"`
	SecondFromAccountID int64 `json:"second_from_account_id"`
	SecondToAccountID   int64 `json:"second_to_account_id"`
	SecondAmount        int64 `json:"second_amount"`
}

// TradeTxResult is the result of the trade transaction
type TradeTxResult struct {
	Trade          Trade    `json:"trade"`
	FirstTransfer  Transfer `json:"first_transfer"`
	SecondTransfer Transfer `json:"second_transfer"`
}

// TradeTx performs a money trade in different currencies.
// It creates the trade and two transfers within a database transaction
func (store *SQLStore) TradeTx(ctx context.Context, arg TradeTxParams) (TradeTxResult, error) {
	var result TradeTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Trade, err = q.CreateTrade(ctx, CreateTradeParams(arg))
		if err != nil {
			return err
		}

		transferResult, err := store.TransferTx(ctx, TransferTxParams{
			FromAccountID: arg.FirstFromAccountID,
			ToAccountID: arg.FirstToAccountID,
			Amount: arg.FirstAmount,
		})
		if err != nil {
			return err
		}
		result.FirstTransfer = transferResult.Transfer

		transferResult, err = store.TransferTx(ctx, TransferTxParams{
			FromAccountID: arg.SecondFromAccountID,
			ToAccountID: arg.SecondToAccountID,
			Amount: arg.SecondAmount,
		})
		if err != nil {
			return err
		}
		result.SecondTransfer = transferResult.Transfer

		return err
	})

	return result, err
}
