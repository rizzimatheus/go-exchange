// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: ask.sql

package db

import (
	"context"
)

const createAsks = `-- name: CreateAsks :one
INSERT INTO asks (pair, from_account_id, to_account_id, price, amount, status) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, pair, from_account_id, to_account_id, price, amount, status, created_at
`

type CreateAsksParams struct {
	Pair          string `json:"pair"`
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Price         int64  `json:"price"`
	Amount        int64  `json:"amount"`
	Status        string `json:"status"`
}

func (q *Queries) CreateAsks(ctx context.Context, arg CreateAsksParams) (Ask, error) {
	row := q.db.QueryRowContext(ctx, createAsks,
		arg.Pair,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Price,
		arg.Amount,
		arg.Status,
	)
	var i Ask
	err := row.Scan(
		&i.ID,
		&i.Pair,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Price,
		&i.Amount,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const getAsk = `-- name: GetAsk :one
SELECT id, pair, from_account_id, to_account_id, price, amount, status, created_at FROM asks
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetAsk(ctx context.Context, id int64) (Ask, error) {
	row := q.db.QueryRowContext(ctx, getAsk, id)
	var i Ask
	err := row.Scan(
		&i.ID,
		&i.Pair,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.Price,
		&i.Amount,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const listAsks = `-- name: ListAsks :many
SELECT id, pair, from_account_id, to_account_id, price, amount, status, created_at FROM asks
WHERE from_account_id = $1 OR to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4
`

type ListAsksParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Limit         int32 `json:"limit"`
	Offset        int32 `json:"offset"`
}

func (q *Queries) ListAsks(ctx context.Context, arg ListAsksParams) ([]Ask, error) {
	rows, err := q.db.QueryContext(ctx, listAsks,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Ask{}
	for rows.Next() {
		var i Ask
		if err := rows.Scan(
			&i.ID,
			&i.Pair,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.Price,
			&i.Amount,
			&i.Status,
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
