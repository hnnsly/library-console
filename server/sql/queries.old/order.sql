-- name: CreateOrder :one
INSERT INTO orders (pharmacy_id, supplier_id, product_id, dosage, quantity, status)
VALUES ($1, $2, $3, $4, $5, $6) -- $6 for status, e.g., 'created'
RETURNING id, pharmacy_id, supplier_id, product_id, dosage, quantity, status, created_at;

-- name: GetOrder :one
SELECT id, pharmacy_id, supplier_id, product_id, dosage, quantity, status, created_at, updated_at
FROM orders
WHERE id = $1;

-- name: ListOrdersByPharmacy :many
SELECT id, supplier_id, product_id, dosage, quantity, status, created_at, updated_at
FROM orders
WHERE pharmacy_id = $1
ORDER BY created_at DESC;

-- name: ListOrdersBySupplier :many
SELECT id, pharmacy_id, product_id, dosage, quantity, status, created_at, updated_at
FROM orders
WHERE supplier_id = $1
ORDER BY created_at DESC;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, pharmacy_id, supplier_id, product_id, dosage, quantity, status, created_at, updated_at;
