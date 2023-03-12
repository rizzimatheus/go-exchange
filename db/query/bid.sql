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

-- name: ListAllBids :many
SELECT id, pair, price, remaining_amount FROM bids
WHERE pair = $1
ORDER BY price DESC
LIMIT $2
OFFSET $3;

-- name: ListTradableBids :many
SELECT * FROM bids
WHERE status = 'active' AND pair = $1 AND price >= $2
ORDER BY price DESC
LIMIT $3
OFFSET $4;

-- name: CreateBid :one
INSERT INTO bids (pair, from_account_id, to_account_id, price, initial_amount, remaining_amount, status) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateBid :one
UPDATE bids
SET 
  status = $2,
  remaining_amount = $3
WHERE id = $1
RETURNING *;
