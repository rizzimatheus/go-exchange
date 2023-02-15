-- name: GetBid :one
SELECT * FROM bids
WHERE id = $1
LIMIT 1;

-- name: ListBids :many
SELECT * FROM bids
WHERE from_account_id = $1 OR to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: CreateBids :one
INSERT INTO bids (pair, from_account_id, to_account_id, price, amount, status) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
