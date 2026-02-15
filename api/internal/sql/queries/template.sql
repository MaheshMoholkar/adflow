-- name: GetTemplateByUserID :many
SELECT * FROM templates WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetTemplateByID :one
SELECT * FROM templates WHERE id = $1 AND user_id = $2;

-- name: CreateTemplate :one
INSERT INTO templates (user_id, name, body, type, channel, image_url, image_key, language, is_default)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateTemplate :one
UPDATE templates
SET name = $3,
    body = $4,
    type = $5,
    channel = $6,
    image_url = $7,
    image_key = $8,
    language = $9,
    is_default = $10,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteTemplate :exec
DELETE FROM templates WHERE id = $1 AND user_id = $2;
