-- name: GetInventoryByPharmacy :many
SELECT product_id, dosage, quantity, price, expiration_date
FROM inventory
WHERE pharmacy_id = @pharmacy_id
ORDER BY product_id, dosage;

-- name: GetInventoryItem :one
SELECT id, product_id, dosage, quantity, price, expiration_date, created_at, updated_at
FROM inventory
WHERE pharmacy_id = @pharmacy_id AND product_id = @product_id AND dosage = @dosage;

-- name: UpsertInventoryItem :one
INSERT INTO inventory (pharmacy_id, product_id, dosage, quantity, price, expiration_date)
VALUES (@pharmacy_id, @product_id, @dosage, @quantity, @price, @expiration_date)
ON CONFLICT (pharmacy_id, product_id, dosage)
DO UPDATE SET
    quantity = EXCLUDED.quantity,
    price = EXCLUDED.price,
    expiration_date = EXCLUDED.expiration_date,
    updated_at = NOW()
RETURNING id, product_id, dosage, quantity, price, expiration_date, created_at, updated_at;

-- name: DeleteInventoryItem :exec
DELETE FROM inventory
WHERE pharmacy_id = @pharmacy_id AND product_id = @product_id AND dosage = @dosage;

-- name: GetExpiredInventoryByPharmacy :many
SELECT product_id, dosage, quantity, expiration_date
FROM inventory
WHERE pharmacy_id = @pharmacy_id AND expiration_date < CURRENT_DATE AND quantity > 0 -- Or just < CURRENT_DATE depending on requirements
ORDER BY expiration_date;

-- name: GetTotalInventoryValueByPharmacy :one
SELECT SUM(quantity * price) AS total_value
FROM inventory
WHERE pharmacy_id = @pharmacy_id;
