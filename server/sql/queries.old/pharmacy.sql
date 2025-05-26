-- name: CreatePharmacy :one
INSERT INTO pharmacies (name)
VALUES ($1)
RETURNING id, name, created_at, updated_at;

-- name: GetPharmacy :one
SELECT id, name, created_at, updated_at
FROM pharmacies
WHERE id = $1;

-- name: ListPharmacies :many
SELECT id, name, created_at, updated_at
FROM pharmacies
ORDER BY name;

-- name: UpdatePharmacy :one
UPDATE pharmacies
SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, name, created_at, updated_at;

-- name: DeletePharmacy :exec
DELETE FROM pharmacies
WHERE id = $1;
