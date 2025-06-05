-- name: CreateReader :one
INSERT INTO readers (user_id, ticket_number, full_name, birth_date, phone, education, reading_hall_id)
VALUES (@user_id, @ticket_number, @full_name, @birth_date, @phone, @education, @reading_hall_id)
RETURNING *;

-- name: GetReaderByID :one
SELECT r.*, u.username, u.email, rh.hall_name, rh.specialization
FROM readers r
JOIN users u ON r.user_id = u.id
LEFT JOIN reading_halls rh ON r.reading_hall_id = rh.id
WHERE r.id = @reader_id AND r.is_active = true;

-- name: GetReaderByUserID :one
SELECT r.*, rh.hall_name, rh.specialization
FROM readers r
LEFT JOIN reading_halls rh ON r.reading_hall_id = rh.id
WHERE r.user_id = @user_id AND r.is_active = true;

-- name: GetReaderByTicketNumber :one
SELECT r.*, u.username, u.email, rh.hall_name, rh.specialization
FROM readers r
JOIN users u ON r.user_id = u.id
LEFT JOIN reading_halls rh ON r.reading_hall_id = rh.id
WHERE r.ticket_number = @ticket_number AND r.is_active = true;

-- name: UpdateReader :one
UPDATE readers
SET
    full_name = COALESCE(@full_name, full_name),
    phone = COALESCE(@phone, phone),
    education = COALESCE(@education, education),
    reading_hall_id = COALESCE(@reading_hall_id, reading_hall_id),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @reader_id
RETURNING *;

-- name: DeactivateReader :exec
UPDATE readers
SET is_active = false, updated_at = CURRENT_TIMESTAMP
WHERE id = @reader_id;

-- name: ListReadersByHall :many
SELECT r.*, u.username, u.email
FROM readers r
JOIN users u ON r.user_id = u.id
WHERE r.reading_hall_id = @hall_id AND r.is_active = true
ORDER BY r.full_name;

-- name: ListAllReaders :many
SELECT r.*, u.username, u.email, rh.hall_name
FROM readers r
JOIN users u ON r.user_id = u.id
LEFT JOIN reading_halls rh ON r.reading_hall_id = rh.id
WHERE r.is_active = true
ORDER BY r.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: SearchReadersByName :many
SELECT r.*, u.username, u.email, rh.hall_name
FROM readers r
JOIN users u ON r.user_id = u.id
LEFT JOIN reading_halls rh ON r.reading_hall_id = rh.id
WHERE r.full_name ILIKE '%' || @search_query || '%' AND r.is_active = true
ORDER BY r.full_name;

-- name: CountReaders :one
SELECT COUNT(*) FROM readers WHERE is_active = true;

-- name: CountReadersByHall :one
SELECT COUNT(*) FROM readers WHERE reading_hall_id = @hall_id AND is_active = true;
