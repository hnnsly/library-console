-- name: CreateProduct :one
INSERT INTO products (name, dosages)
VALUES ($1, $2)
RETURNING id, name, dosages, created_at, updated_at;

-- name: GetProduct :one
SELECT id, name, dosages, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, dosages, created_at, updated_at
FROM products
ORDER BY name;
