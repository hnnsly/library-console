-- name: CreateReader :one
INSERT INTO readers (ticket_number, full_name, email, phone)
VALUES (@ticket_number, @full_name, @email, @phone)
RETURNING id, ticket_number, created_at;

-- name: UpdateReader :one
UPDATE readers
SET full_name = @full_name, email = @email, phone = @phone
WHERE id = @id
RETURNING id, ticket_number, full_name, email, phone;

-- name: DeactivateReader :exec
UPDATE readers
SET is_active = false
WHERE id = @id;

-- name: GetReaderByTicketNumber :one
SELECT id, ticket_number, full_name, email, phone, is_active, registration_date
FROM readers
WHERE ticket_number = @ticket_number;

-- name: GetReaderById :one
SELECT id, ticket_number, full_name, email, phone, is_active, registration_date
FROM readers
WHERE id = @id;

-- name: SearchReaders :many
SELECT id, ticket_number, full_name, email, phone, is_active, registration_date
FROM readers
WHERE (full_name ILIKE '%' || @search_term || '%'
    OR ticket_number ILIKE '%' || @search_term || '%')
  AND (@include_inactive::boolean OR is_active = true)
ORDER BY full_name;

-- name: GetActiveReaders :many
SELECT id, ticket_number, full_name, email, phone, registration_date
FROM readers
WHERE is_active = true
ORDER BY full_name;

-- name: CheckReaderOverdueBooks :one
SELECT COUNT(*) as overdue_books
FROM book_issues bi
WHERE bi.reader_id = @reader_id
  AND bi.return_date IS NULL
  AND bi.due_date < CURRENT_DATE;
