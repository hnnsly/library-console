-- name: CreateReadingHall :one
INSERT INTO reading_halls (library_name, hall_name, specialization, total_seats)
VALUES (@library_name, @hall_name, @specialization, @total_seats)
RETURNING *;

-- name: GetReadingHallByID :one
SELECT * FROM reading_halls WHERE id = @hall_id;

-- name: UpdateReadingHall :one
UPDATE reading_halls
SET
    library_name = COALESCE(@library_name, library_name),
    hall_name = COALESCE(@hall_name, hall_name),
    specialization = COALESCE(@specialization, specialization),
    total_seats = COALESCE(@total_seats, total_seats),
    occupied_seats = COALESCE(@occupied_seats, occupied_seats)
WHERE id = @hall_id
RETURNING *;

-- name: DeleteReadingHall :exec
DELETE FROM reading_halls WHERE id = @hall_id;

-- name: ListReadingHalls :many
SELECT * FROM reading_halls ORDER BY library_name, hall_name;

-- name: GetHallStatistics :many
SELECT
    rh.id,
    rh.library_name,
    rh.hall_name,
    rh.specialization,
    rh.total_seats,
    rh.occupied_seats,
    (rh.total_seats - rh.occupied_seats) as free_seats,
    COUNT(r.id) as registered_readers,
    COUNT(CASE WHEN r.is_active THEN 1 END) as active_readers
FROM reading_halls rh
LEFT JOIN readers r ON rh.id = r.reading_hall_id
GROUP BY rh.id, rh.library_name, rh.hall_name, rh.specialization, rh.total_seats, rh.occupied_seats
ORDER BY rh.library_name, rh.hall_name;

-- name: UpdateHallOccupancy :exec
UPDATE reading_halls
SET occupied_seats = @occupied_seats
WHERE id = @hall_id;
