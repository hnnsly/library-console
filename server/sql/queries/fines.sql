-- name: CreateFine :one
INSERT INTO fines (reader_id, book_issue_id, amount, reason, fine_date, librarian_id)
VALUES (@reader_id, @book_issue_id, @amount, @reason, @fine_date, @librarian_id)
RETURNING *;

-- name: GetFineByID :one
SELECT
    f.*,
    r.full_name as reader_name,
    r.ticket_number,
    u.username as librarian_name
FROM fines f
JOIN readers r ON f.reader_id = r.id
LEFT JOIN users u ON f.librarian_id = u.id
WHERE f.id = @fine_id;

-- name: UpdateFine :one
UPDATE fines
SET
    amount = COALESCE(@amount, amount),
    reason = COALESCE(@reason, reason),
    paid_date = COALESCE(@paid_date, paid_date),
    paid_amount = COALESCE(@paid_amount, paid_amount),
    is_paid = COALESCE(@is_paid, is_paid),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @fine_id
RETURNING *;

-- name: PayFine :exec
UPDATE fines
SET
    paid_date = @paid_date,
    paid_amount = @paid_amount,
    is_paid = (@paid_amount >= amount),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @fine_id;

-- name: DeleteFine :exec
DELETE FROM fines WHERE id = @fine_id;

-- name: GetFinesByReader :many
SELECT f.*, bi.issue_date, bi.due_date, b.title as book_title
FROM fines f
LEFT JOIN book_issues bi ON f.book_issue_id = bi.id
LEFT JOIN book_copies bc ON bi.book_copy_id = bc.id
LEFT JOIN books b ON bc.book_id = b.id
WHERE f.reader_id = @reader_id
ORDER BY f.created_at DESC;

-- name: GetUnpaidFinesByReader :many
SELECT f.*, bi.issue_date, bi.due_date, b.title as book_title
FROM fines f
LEFT JOIN book_issues bi ON f.book_issue_id = bi.id
LEFT JOIN book_copies bc ON bi.book_copy_id = bc.id
LEFT JOIN books b ON bc.book_id = b.id
WHERE f.reader_id = @reader_id AND f.is_paid = false
ORDER BY f.fine_date;

-- name: GetAllUnpaidFines :many
SELECT
    f.*,
    r.full_name as reader_name,
    r.ticket_number,
    b.title as book_title
FROM fines f
JOIN readers r ON f.reader_id = r.id
LEFT JOIN book_issues bi ON f.book_issue_id = bi.id
LEFT JOIN book_copies bc ON bi.book_copy_id = bc.id
LEFT JOIN books b ON bc.book_id = b.id
WHERE f.is_paid = false
ORDER BY f.fine_date;

-- name: GetTotalDebtByReader :one
SELECT COALESCE(SUM(amount - paid_amount), 0) as total_debt
FROM fines
WHERE reader_id = @reader_id AND is_paid = false;

-- name: GetFineStatistics :one
SELECT
    COUNT(*) as total_fines,
    COUNT(CASE WHEN is_paid THEN 1 END) as paid_fines,
    COUNT(CASE WHEN NOT is_paid THEN 1 END) as unpaid_fines,
    COALESCE(SUM(amount), 0) as total_amount,
    COALESCE(SUM(paid_amount), 0) as total_paid,
    COALESCE(SUM(CASE WHEN NOT is_paid THEN amount - paid_amount ELSE 0 END), 0) as total_debt
FROM fines
WHERE fine_date >= @from_date AND fine_date <= @to_date;
