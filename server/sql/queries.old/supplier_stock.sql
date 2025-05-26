-- name: GetStockBySupplier :many
SELECT product_id, dosage, quantity, preference_level
FROM supplier_stock
WHERE supplier_id = $1
ORDER BY preference_level DESC, product_id, dosage;

-- name: GetSupplierStockItem :one
SELECT id, product_id, dosage, quantity, preference_level, created_at, updated_at
FROM supplier_stock
WHERE supplier_id = $1 AND product_id = $2 AND dosage = $3;

-- name: UpsertSupplierStockItem :one
INSERT INTO supplier_stock (supplier_id, product_id, dosage, quantity, preference_level)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (supplier_id, product_id, dosage)
DO UPDATE SET
    quantity = EXCLUDED.quantity,
    preference_level = EXCLUDED.preference_level,
    updated_at = NOW()
RETURNING id, product_id, dosage, quantity, preference_level, created_at, updated_at;

-- name: DeleteSupplierStockItem :exec
DELETE FROM supplier_stock
WHERE supplier_id = $1 AND product_id = $2 AND dosage = $3;

-- name: SearchSupplierStockByProductAndDosage :many
SELECT
    s.id AS supplier_id,
    s.name AS supplier_name,
    ss.quantity AS available_quantity,
    ss.preference_level
FROM supplier_stock ss
JOIN suppliers s ON ss.supplier_id = s.id
JOIN products p ON ss.product_id = p.id
WHERE p.name ILIKE $1 AND ss.dosage = $2 AND ss.quantity > 0 -- Assuming productName is a partial match
ORDER BY ss.preference_level DESC, s.name;
