-- name: GetAllHalls :many
SELECT id, name, specialization, total_seats, occupied_seats,
       (total_seats - occupied_seats) as free_seats,
       ROUND((occupied_seats * 100.0 / total_seats), 2) as occupancy_percent,
       working_hours, status
FROM halls
WHERE status = 'open'
ORDER BY name;

-- name: GetHallByID :one
SELECT * FROM halls WHERE id = @hall_id;

-- name: UpdateHallOccupancy :exec
UPDATE halls
SET occupied_seats = (
    SELECT COUNT(DISTINCT reader_id)
    FROM loan_history lh
    JOIN readers r ON lh.reader_id = r.id
    WHERE r.hall_id = @hall_id AND lh.status = 'active'
),
average_occupancy = (occupied_seats * 100.0 / total_seats),
updated_at = NOW()
WHERE id = @hall_id;

-- name: GetHallStatistics :many
SELECT
    h.name as hall_name,
    h.specialization,
    COUNT(DISTINCT b.id) as total_books,
    SUM(b.total_copies) as total_copies,
    COUNT(lh.id) as total_loans,
    ROUND(h.average_occupancy, 2) as avg_occupancy_percent,
    (h.total_seats - h.occupied_seats) as current_free_seats
FROM halls h
LEFT JOIN books b ON h.id = b.hall_id
LEFT JOIN loan_history lh ON b.id = lh.book_id AND lh.loan_date >= CURRENT_DATE - INTERVAL '@days_back days'
GROUP BY h.id, h.name, h.specialization, h.average_occupancy, h.total_seats, h.occupied_seats
ORDER BY total_loans DESC;
