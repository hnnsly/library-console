-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, role)
VALUES (@username, @email, @password_hash, @role)
RETURNING id, username, email, role, is_active, created_at;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, role, is_active
FROM users
WHERE username = @username;

-- name: GetUserById :one
SELECT id, username, email, role, is_active, created_at
FROM users
WHERE id = @id;

-- name: UpdateUser :one
UPDATE users
SET email = @email, role = @role
WHERE id = @id
RETURNING id, username, email, role, is_active;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = false
WHERE id = @id;

-- name: GetAllUsers :many
SELECT id, username, email, role, is_active, created_at
FROM users
WHERE is_active = true
ORDER BY username;
