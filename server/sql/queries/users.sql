-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, role)
VALUES (@username, @email, @password_hash, @role)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = @user_id AND is_active = true;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = @username AND is_active = true;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = @email AND is_active = true;

-- name: UpdateUser :one
UPDATE users
SET
    username = COALESCE(@username, username),
    email = COALESCE(@email, email),
    role = COALESCE(@role, role),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @user_id
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = @password_hash, updated_at = CURRENT_TIMESTAMP
WHERE id = @user_id;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = false, updated_at = CURRENT_TIMESTAMP
WHERE id = @user_id;

-- name: ListUsers :many
SELECT * FROM users
WHERE is_active = true
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: ListUsersByRole :many
SELECT * FROM users
WHERE role = @role AND is_active = true
ORDER BY created_at DESC;

-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE is_active = true;

-- name: CountUsersByRole :one
SELECT COUNT(*) FROM users WHERE role = @role AND is_active = true;
