-- name: CreateEntry :one
INSERT INTO entries (account_id, amount) VALUES ($1, $2) RETURNING *;

-- name: GetListEntry :many
SELECT id, account_id, amount, created_at FROM entries;

-- name: GetByIDEntry :one
SELECT id, account_id, amount, created_at FROM entries WHERE id = $1;

