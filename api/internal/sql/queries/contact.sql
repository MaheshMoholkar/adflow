-- name: GetContactsByUserID :many
SELECT * FROM contacts WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpsertContact :one
INSERT INTO contacts (user_id, phone, name)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, phone) DO UPDATE
SET name = EXCLUDED.name
RETURNING *;

-- name: UpsertContactBatch :exec
INSERT INTO contacts (user_id, phone, name)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, phone) DO UPDATE
SET name = EXCLUDED.name;
