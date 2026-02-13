-- name: CreateUser :one
INSERT INTO users (phone, phone_verified, password_hash, name, business_name, city, address)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByPhone :one
SELECT * FROM users WHERE phone = $1;

-- name: UpdateUser :one
UPDATE users
SET name = COALESCE($2, name),
    business_name = COALESCE($3, business_name),
    city = COALESCE($4, city),
    address = COALESCE($5, address),
    location_url = COALESCE($6, location_url),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPlan :exec
UPDATE users
SET plan = $2,
    plan_started_at = NOW(),
    plan_expires_at = '2099-12-31T23:59:59Z',
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserStatus :exec
UPDATE users SET status = $2, updated_at = NOW() WHERE id = $1;

-- name: ListAllUsers :many
SELECT * FROM users ORDER BY created_at DESC;
