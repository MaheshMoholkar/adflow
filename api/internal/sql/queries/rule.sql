-- name: GetRuleByUserID :one
SELECT * FROM rules WHERE user_id = $1;

-- name: UpsertRule :one
INSERT INTO rules (user_id, config, updated_at)
VALUES ($1, $2, NOW())
ON CONFLICT (user_id) DO UPDATE
SET config = $2, updated_at = NOW()
RETURNING *;
