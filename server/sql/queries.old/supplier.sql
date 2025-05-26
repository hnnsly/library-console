-- name: CreateSupplier :one
INSERT INTO suppliers (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at;

-- name: GetSupplier :one
SELECT id, name, created_at, updated_at
FROM suppliers
WHERE id = $1;

-- name: ListSuppliers :many
SELECT id, name, created_at, updated_at
FROM suppliers
ORDER BY name;
