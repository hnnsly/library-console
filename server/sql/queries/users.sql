-- name: CreateUser :one
INSERT INTO users (
    username, email, password_hash, role, full_name, phone, created_by
) VALUES (
    @username, @email, @password_hash, @role, @full_name, @phone, @created_by
) RETURNING id, username, email, role, full_name, phone, is_active, created_at, updated_at;

-- name: GetUserByUsername :one
SELECT id, username, email, password_hash, role, full_name, phone, is_active, is_first_admin, last_login_at, created_at, updated_at
FROM users
WHERE username = @username AND is_active = true;

-- name: GetUserByID :one
SELECT id, username, email, role, full_name, phone, is_active, is_first_admin, last_login_at, created_at, updated_at
FROM users
WHERE id = @id;

-- name: GetAllUsers :many
SELECT id, username, email, role, full_name, phone, is_active, is_first_admin, last_login_at, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT @limit_users OFFSET @offset_users;

-- name: UpdateUser :one
UPDATE users
SET
    email = @email,
    full_name = @full_name,
    phone = @phone,
    updated_at = NOW()
WHERE id = @id
RETURNING id, username, email, role, full_name, phone, is_active, created_at, updated_at;

-- name: UpdateUserRole :exec
UPDATE users
SET role = @role, updated_at = NOW()
WHERE id = @id AND is_first_admin = false;

-- name: DeactivateUser :exec
UPDATE users
SET is_active = false, updated_at = NOW()
WHERE id = @id AND is_first_admin = false;

-- name: ActivateUser :exec
UPDATE users
SET is_active = true, updated_at = NOW()
WHERE id = @id;

-- name: UpdatePassword :exec
UPDATE users
SET password_hash = @password_hash, updated_at = NOW()
WHERE id = @id;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login_at = NOW()
WHERE id = @id;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = @id AND is_first_admin = false;

-- name: GetUsersByRole :many
SELECT id, username, email, role, full_name, phone, is_active, last_login_at, created_at
FROM users
WHERE role = @role AND is_active = true
ORDER BY full_name;
