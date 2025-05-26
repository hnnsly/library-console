-- name: CreateFine :one
INSERT INTO fines (
    loan_history_id, reader_id, fine_type, amount, fine_date, description, librarian_id
) VALUES (
    @loan_history_id, @reader_id, @fine_type, @amount, CURRENT_DATE, @description, @librarian_id
) RETURNING *;

-- name: PayFine :exec
UPDATE fines
SET status = 'paid', payment_date = CURRENT_DATE, updated_at = NOW()
WHERE id = @fine_id;

-- name: WaiveFine :exec
UPDATE fines
SET status = 'waived', updated_at = NOW()
WHERE id = @fine_id;

-- name: GetReaderFines :many
SELECT f.*, lh.loan_date, lh.due_date, b.title, b.author
FROM fines f
JOIN loan_history lh ON f.loan_history_id = lh.id
JOIN books b ON lh.book_id = b.id
WHERE f.reader_id = @reader_id
ORDER BY f.fine_date DESC;

-- name: GetUnpaidFines :many
SELECT
    f.*,
    r.full_name as reader_name,
    r.ticket_number,
    r.phone,
    b.title as book_title
FROM fines f
JOIN readers r ON f.reader_id = r.id
JOIN loan_history lh ON f.loan_history_id = lh.id
JOIN books b ON lh.book_id = b.id
WHERE f.status = 'unpaid'
ORDER BY f.fine_date DESC;

-- name: GetDebtorReaders :many
SELECT
    r.full_name,
    r.ticket_number,
    r.phone,
    r.email,
    SUM(f.amount) as total_debt,
    COUNT(f.id) as fine_count
FROM readers r
JOIN fines f ON r.id = f.reader_id
WHERE f.status = 'unpaid'
GROUP BY r.id, r.full_name, r.ticket_number, r.phone, r.email
ORDER BY total_debt DESC;

-- name: CalculateOverdueFine :one
SELECT
    id as loan_history_id,
    GREATEST(0, EXTRACT(DAYS FROM (CURRENT_DATE - due_date))::int) as overdue_days,
    GREATEST(0, EXTRACT(DAYS FROM (CURRENT_DATE - due_date))::int) * @daily_fine_rate as calculated_fine
FROM loan_history
WHERE id = @loan_history_id;
