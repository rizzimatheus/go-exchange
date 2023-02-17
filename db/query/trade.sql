-- name: GetTrade :one
SELECT * FROM trades
WHERE id = $1
LIMIT 1;

-- name: ListTrades :many
SELECT * FROM trades
WHERE first_from_account_id = $1 OR first_to_account_id = $2 OR second_from_account_id = $3 OR second_to_account_id = $4
ORDER BY id
LIMIT $5
OFFSET $6;

-- name: CreateTrade :one
INSERT INTO trades (first_from_account_id, first_to_account_id, first_amount, second_from_account_id, second_to_account_id, second_amount) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
