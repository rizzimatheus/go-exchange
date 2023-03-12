-- name: GetAsk :one
SELECT * FROM asks
WHERE id = $1
LIMIT 1;

-- name: ListAsks :many
SELECT * FROM asks
WHERE from_account_id = $1 OR to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;

-- name: ListAllAsks :many
SELECT id, pair, price, initial_amount, remaining_amount FROM asks
WHERE pair = $1
ORDER BY price ASC
LIMIT $2
OFFSET $3;

-- name: ListTradableAsks :many
SELECT * FROM asks
WHERE status = 'active' AND pair = $1 AND price <= $2
ORDER BY price DESC
LIMIT $3
OFFSET $4;

-- name: CreateAsk :one
INSERT INTO asks (pair, from_account_id, to_account_id, price, initial_amount, remaining_amount, status) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateAsk :one
UPDATE asks
SET 
  status = $2,
  remaining_amount = $3
WHERE id = $1
RETURNING *;
