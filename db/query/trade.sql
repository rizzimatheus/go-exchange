-- name: GetTrade :one
SELECT * FROM trades
WHERE id = $1
LIMIT 1;

-- name: ListTrades :many
SELECT * FROM trades
WHERE first_transfer_id = $1 OR second_transfer_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: CreateTrade :one
INSERT INTO trades (first_transfer_id, second_transfer_id) VALUES ($1, $2)
RETURNING *;
