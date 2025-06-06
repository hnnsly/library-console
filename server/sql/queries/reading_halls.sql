-- name: CreateReadingHall :one
INSERT INTO reading_halls (hall_name, specialization, total_seats)
VALUES (@hall_name, @specialization, @total_seats)
RETURNING id, hall_name, specialization, total_seats, current_visitors;

-- name: UpdateReadingHall :one
UPDATE reading_halls
SET hall_name = @hall_name, specialization = @specialization, total_seats = @total_seats
WHERE id = @id
RETURNING id, hall_name, specialization, total_seats, current_visitors;

-- name: GetReadingHallById :one
SELECT id, hall_name, specialization, total_seats, current_visitors
FROM reading_halls
WHERE id = @id;

-- name: GetAllReadingHalls :many
SELECT id, hall_name, specialization, total_seats, current_visitors
FROM reading_halls
ORDER BY hall_name;

-- name: GetHallsDashboard :many
SELECT
    rh.id,
    rh.hall_name,
    rh.specialization,
    rh.total_seats,
    rh.current_visitors,
    ROUND(
        (rh.current_visitors::numeric / rh.total_seats * 100), 2
    ) as occupancy_percentage,
    (rh.total_seats - rh.current_visitors) as free_seats,
    COALESCE(daily_stats.visits_today, 0) as visits_today,
    COALESCE(daily_stats.unique_visitors_today, 0) as unique_visitors_today
FROM reading_halls rh
LEFT JOIN (
    SELECT
        hall_id,
        COUNT(*) as visits_today,
        COUNT(DISTINCT reader_id) as unique_visitors_today
    FROM hall_visits
    WHERE DATE(visit_time) = CURRENT_DATE
    GROUP BY hall_id
) daily_stats ON rh.id = daily_stats.hall_id
ORDER BY rh.hall_name;

-- name: UpdateHallVisitorCount :exec
UPDATE reading_halls
SET current_visitors = current_visitors + @change
WHERE id = @hall_id;
