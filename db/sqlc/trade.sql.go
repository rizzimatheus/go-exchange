// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: trade.sql

package db

import (
	"context"
)

const createTrade = `-- name: CreateTrade :one
INSERT INTO trades (first_from_account_id, first_to_account_id, first_amount, second_from_account_id, second_to_account_id, second_amount) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, first_from_account_id, first_to_account_id, first_amount, second_from_account_id, second_to_account_id, second_amount, created_at
`

type CreateTradeParams struct {
	FirstFromAccountID  int64 `json:"first_from_account_id"`
	FirstToAccountID    int64 `json:"first_to_account_id"`
	FirstAmount         int64 `json:"first_amount"`
	SecondFromAccountID int64 `json:"second_from_account_id"`
	SecondToAccountID   int64 `json:"second_to_account_id"`
	SecondAmount        int64 `json:"second_amount"`
}

func (q *Queries) CreateTrade(ctx context.Context, arg CreateTradeParams) (Trade, error) {
	row := q.db.QueryRowContext(ctx, createTrade,
		arg.FirstFromAccountID,
		arg.FirstToAccountID,
		arg.FirstAmount,
		arg.SecondFromAccountID,
		arg.SecondToAccountID,
		arg.SecondAmount,
	)
	var i Trade
	err := row.Scan(
		&i.ID,
		&i.FirstFromAccountID,
		&i.FirstToAccountID,
		&i.FirstAmount,
		&i.SecondFromAccountID,
		&i.SecondToAccountID,
		&i.SecondAmount,
		&i.CreatedAt,
	)
	return i, err
}

const getTrade = `-- name: GetTrade :one
SELECT id, first_from_account_id, first_to_account_id, first_amount, second_from_account_id, second_to_account_id, second_amount, created_at FROM trades
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetTrade(ctx context.Context, id int64) (Trade, error) {
	row := q.db.QueryRowContext(ctx, getTrade, id)
	var i Trade
	err := row.Scan(
		&i.ID,
		&i.FirstFromAccountID,
		&i.FirstToAccountID,
		&i.FirstAmount,
		&i.SecondFromAccountID,
		&i.SecondToAccountID,
		&i.SecondAmount,
		&i.CreatedAt,
	)
	return i, err
}

const listTrades = `-- name: ListTrades :many
SELECT id, first_from_account_id, first_to_account_id, first_amount, second_from_account_id, second_to_account_id, second_amount, created_at FROM trades
WHERE first_from_account_id = $1 OR first_to_account_id = $2 OR second_from_account_id = $3 OR second_to_account_id = $4
ORDER BY id
LIMIT $5
OFFSET $6
`

type ListTradesParams struct {
	FirstFromAccountID  int64 `json:"first_from_account_id"`
	FirstToAccountID    int64 `json:"first_to_account_id"`
	SecondFromAccountID int64 `json:"second_from_account_id"`
	SecondToAccountID   int64 `json:"second_to_account_id"`
	Limit               int32 `json:"limit"`
	Offset              int32 `json:"offset"`
}

func (q *Queries) ListTrades(ctx context.Context, arg ListTradesParams) ([]Trade, error) {
	rows, err := q.db.QueryContext(ctx, listTrades,
		arg.FirstFromAccountID,
		arg.FirstToAccountID,
		arg.SecondFromAccountID,
		arg.SecondToAccountID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Trade{}
	for rows.Next() {
		var i Trade
		if err := rows.Scan(
			&i.ID,
			&i.FirstFromAccountID,
			&i.FirstToAccountID,
			&i.FirstAmount,
			&i.SecondFromAccountID,
			&i.SecondToAccountID,
			&i.SecondAmount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
