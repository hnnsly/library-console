-- name: CreateLoan :one
INSERT INTO loan_history (
    book_id, reader_id, librarian_id, loan_date, due_date, status
) VALUES (
    @book_id, @reader_id, @librarian_id, CURRENT_DATE,
    CURRENT_DATE + INTERVAL '@loan_days days', 'active'
) RETURNING *;

-- name: CheckLoanEligibility :one
SELECT
    b.available_copies > 0 as book_available,
    (SELECT COUNT(*) FROM loan_history lh WHERE lh.reader_id = @reader_id AND lh.status = 'active') as current_loans,
    r.max_books_allowed,
    r.status as reader_status,
    r.total_debt,
    (SELECT COUNT(*) FROM reservations res WHERE res.book_id = @book_id AND res.status = 'active' AND res.reader_id != @reader_id) as queue_length
FROM books b, readers r
WHERE b.id = @book_id AND r.id = @reader_id;

-- name: ReturnBook :exec
UPDATE loan_history
SET status = 'returned',
    return_date = CURRENT_DATE,
    return_librarian_id = @librarian_id,
    updated_at = NOW()
WHERE id = @loan_id AND status = 'active';

-- name: RenewLoan :exec
UPDATE loan_history
SET due_date = due_date + INTERVAL '@extension_days days',
    renewals_count = renewals_count + 1,
    updated_at = NOW()
WHERE id = @loan_id AND status = 'active';

-- name: CreateRenewal :exec
INSERT INTO renewals (
    loan_history_id, renewal_date, old_due_date, new_due_date, librarian_id, reason
) VALUES (
    @loan_history_id, CURRENT_DATE, @old_due_date, @new_due_date, @librarian_id, @reason
);

-- name: GetReaderCurrentLoans :many
SELECT b.title, b.author, b.book_code, lh.loan_date, lh.due_date,
       lh.renewals_count,
       EXTRACT(DAYS FROM (CURRENT_DATE - lh.due_date))::int as days_overdue,
       CASE WHEN lh.due_date < CURRENT_DATE THEN 'Просрочена' ELSE 'В срок' END as status_text
FROM loan_history lh
JOIN books b ON lh.book_id = b.id
WHERE lh.reader_id = @reader_id AND lh.status = 'active'
ORDER BY lh.due_date;

-- name: GetReaderLoanHistory :many
SELECT b.title, b.author, lh.loan_date, lh.return_date, lh.due_date,
       lh.status, lh.renewals_count,
       CASE WHEN lh.return_date > lh.due_date
            THEN EXTRACT(DAYS FROM (lh.return_date - lh.due_date))::int
            ELSE 0 END as overdue_days
FROM loan_history lh
JOIN books b ON lh.book_id = b.id
WHERE lh.reader_id = @reader_id
ORDER BY lh.loan_date DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: GetOverdueBooks :many
SELECT
    lh.id,
    b.title,
    b.author,
    b.book_code,
    r.full_name as reader_name,
    r.ticket_number,
    r.phone,
    r.email,
    lh.loan_date,
    lh.due_date,
    EXTRACT(DAYS FROM (CURRENT_DATE - lh.due_date))::int as days_overdue,
    CASE
        WHEN EXTRACT(DAYS FROM (CURRENT_DATE - lh.due_date)) <= 7 THEN 'Легкая'
        WHEN EXTRACT(DAYS FROM (CURRENT_DATE - lh.due_date)) <= 30 THEN 'Средняя'
        ELSE 'Критическая'
    END as overdue_level
FROM loan_history lh
JOIN books b ON lh.book_id = b.id
JOIN readers r ON lh.reader_id = r.id
WHERE lh.status = 'active' AND lh.due_date < CURRENT_DATE
ORDER BY days_overdue DESC
LIMIT @result_limit;

-- name: GetBooksDueToday :many
SELECT
    b.title,
    b.author,
    r.full_name as reader_name,
    r.phone,
    lh.loan_date,
    lh.renewals_count
FROM loan_history lh
JOIN books b ON lh.book_id = b.id
JOIN readers r ON lh.reader_id = r.id
WHERE lh.status = 'active' AND lh.due_date = CURRENT_DATE
ORDER BY r.full_name;

-- name: GetActiveLoansByBook :many
SELECT
    lh.id,
    r.full_name as reader_name,
    r.ticket_number,
    lh.loan_date,
    lh.due_date,
    lh.renewals_count
FROM loan_history lh
JOIN readers r ON lh.reader_id = r.id
WHERE lh.book_id = @book_id AND lh.status = 'active'
ORDER BY lh.loan_date DESC;

-- name: GetLoanByID :one
SELECT lh.*, b.title, b.author, r.full_name as reader_name
FROM loan_history lh
JOIN books b ON lh.book_id = b.id
JOIN readers r ON lh.reader_id = r.id
WHERE lh.id = @loan_id;

-- name: MarkLoanAsLost :exec
UPDATE loan_history
SET
    status = 'lost',
    updated_at = NOW()
WHERE id = $1;
