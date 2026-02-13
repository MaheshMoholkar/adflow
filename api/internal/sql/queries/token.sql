-- name: CreateToken :one
INSERT INTO tokens (user_id, token, expires_at, token_type, client_ip, user_agent)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTokenByToken :one
SELECT * FROM tokens WHERE token = $1;

-- name: RevokeToken :exec
UPDATE tokens SET is_revoked = true WHERE token = $1;

-- name: RevokeAllUserTokens :exec
UPDATE tokens SET is_revoked = true WHERE user_id = $1;

-- name: RevokeAllUserTokensByType :exec
UPDATE tokens SET is_revoked = true WHERE user_id = $1 AND token_type = $2;

-- name: UpdateTokenLastUsed :exec
UPDATE tokens SET last_used_at = NOW() WHERE id = $1;

-- name: DeleteExpiredTokens :exec
DELETE FROM tokens WHERE expires_at < $1;
