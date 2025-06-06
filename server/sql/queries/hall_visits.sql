-- name: RegisterHallEntry :one
INSERT INTO hall_visits (reader_id, hall_id, visit_type, librarian_id)
VALUES (
    (SELECT id FROM readers WHERE ticket_number = @ticket_number),
    @hall_id,
    'entry',
    @librarian_id
)
RETURNING id, visit_time;

-- name: RegisterHallExit :one
INSERT INTO hall_visits (reader_id, hall_id, visit_type, librarian_id)
VALUES (
    (SELECT id FROM readers WHERE ticket_number = @ticket_number),
    @hall_id,
    'exit',
    @librarian_id
)
RETURNING id, visit_time;

-- name: GetHourlyVisitStats :many
SELECT
    EXTRACT(HOUR FROM visit_time) as hour,
    COUNT(*) as visits_count,
    COUNT(DISTINCT reader_id) as unique_visitors
FROM hall_visits
WHERE hall_id = @hall_id
  AND DATE(visit_time) = @visit_date
  AND visit_type = 'entry'
GROUP BY EXTRACT(HOUR FROM visit_time)
ORDER BY hour;

-- name: GetDailyVisitStats :many
SELECT
    DATE(visit_time) as visit_date,
    COUNT(*) as total_visits,
    COUNT(DISTINCT reader_id) as unique_visitors
FROM hall_visits
WHERE hall_id = @hall_id
  AND visit_time >= @start_date
  AND visit_time <= @end_date
  AND visit_type = 'entry'
GROUP BY DATE(visit_time)
ORDER BY visit_date;

-- name: GetRecentHallVisits :many
SELECT
    hv.visit_time,
    hv.visit_type,
    r.full_name as reader_name,
    r.ticket_number,
    rh.hall_name,
    u.username as librarian_name
FROM hall_visits hv
JOIN readers r ON hv.reader_id = r.id
JOIN reading_halls rh ON hv.hall_id = rh.id
LEFT JOIN users u ON hv.librarian_id = u.id
WHERE hv.visit_time >= @since_date
ORDER BY hv.visit_time DESC
LIMIT @limit_count;

-- name: GetReaderVisitHistory :many
SELECT
    hv.visit_time,
    hv.visit_type,
    rh.hall_name,
    u.username as librarian_name
FROM hall_visits hv
JOIN reading_halls rh ON hv.hall_id = rh.id
LEFT JOIN users u ON hv.librarian_id = u.id
WHERE hv.reader_id = @reader_id
ORDER BY hv.visit_time DESC;
