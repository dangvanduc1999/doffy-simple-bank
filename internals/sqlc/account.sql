-- name: Create :one
INSERT INTO accounts (owner, balance, currency) VALUES ($1, $2, $3) RETURNING *;

-- name: GetList :many
SELECT id, owner, balance, currency, created_at FROM accounts;

-- name: GetByID :one
SELECT id, owner, balance, currency, created_at FROM accounts WHERE id = $1;

-- name: UpdateOwner :exec
UPDATE accounts SET owner = $2 WHERE id = $1;

-- name: DeleteByID :exec
DELETE FROM accounts WHERE id = $1;