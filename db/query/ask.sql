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

-- name: CreateAsk :one
INSERT INTO asks (pair, from_account_id, to_account_id, price, amount, status) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateAsk :one
UPDATE asks
  SET status = $2
WHERE id = $1
RETURNING *;
