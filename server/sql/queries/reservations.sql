-- name: CreateReservation :one
INSERT INTO reservations (
    book_id, reader_id, reservation_date, expiration_date, priority_order
) VALUES (
    @book_id, @reader_id, CURRENT_DATE, CURRENT_DATE + INTERVAL '7 days',
    (SELECT COALESCE(MAX(priority_order), 0) + 1 FROM reservations WHERE book_id = @book_id AND status = 'active')
) RETURNING *;

-- name: GetBookQueue :many
SELECT
    r.full_name,
    r.ticket_number,
    res.reservation_date,
    res.priority_order,
    res.expiration_date
FROM reservations res
JOIN readers r ON res.reader_id = r.id
WHERE res.book_id = @book_id AND res.status = 'active'
ORDER BY res.priority_order;

-- name: FulfillReservation :exec
UPDATE reservations
SET status = 'fulfilled', updated_at = NOW()
WHERE id = @reservation_id;

-- name: CancelReservation :exec
UPDATE reservations
SET status = 'cancelled', updated_at = NOW()
WHERE id = @reservation_id;

-- name: GetReaderReservations :many
SELECT res.*, b.title, b.author, b.book_code
FROM reservations res
JOIN books b ON res.book_id = b.id
WHERE res.reader_id = @reader_id AND res.status = 'active'
ORDER BY res.reservation_date DESC;

-- name: GetExpiredReservations :many
SELECT res.*, r.full_name, b.title
FROM reservations res
JOIN readers r ON res.reader_id = r.id
JOIN books b ON res.book_id = b.id
WHERE res.status = 'active' AND res.expiration_date < CURRENT_DATE;
