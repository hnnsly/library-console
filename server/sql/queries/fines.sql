-- name: CreateFine :one
INSERT INTO fines (reader_id, book_issue_id, amount, reason)
VALUES (@reader_id, @book_issue_id, @amount, @reason)
RETURNING id, fine_date, amount, reason;

-- name: PayFine :one
UPDATE fines
SET paid_date = CURRENT_DATE, is_paid = true
WHERE id = @fine_id
RETURNING id, amount, paid_date;

-- name: GetReaderFines :many
SELECT id, amount, reason, fine_date, paid_date, is_paid
FROM fines
WHERE reader_id = @reader_id
ORDER BY fine_date DESC;

-- name: GetUnpaidFines :many
SELECT
    f.id,
    r.ticket_number,
    r.full_name as reader_name,
    f.amount,
    f.reason,
    f.fine_date
FROM fines f
JOIN readers r ON f.reader_id = r.id
WHERE f.is_paid = false
ORDER BY f.fine_date;

-- name: GetReaderUnpaidFinesTotal :one
SELECT COALESCE(SUM(amount), 0) as total_unpaid
FROM fines
WHERE reader_id = @reader_id AND is_paid = false;
