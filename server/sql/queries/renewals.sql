-- name: CreateRenewalRecord :one
INSERT INTO renewals (
    loan_history_id, renewal_date, old_due_date, new_due_date, librarian_id, reason
) VALUES (
    @loan_history_id, CURRENT_DATE, @old_due_date, @new_due_date, @librarian_id, @reason
) RETURNING *;

-- name: GetRenewalsForLoan :many
SELECT
    r.*,
    lib.full_name as librarian_name
FROM renewals r
JOIN librarians lib ON r.librarian_id = lib.id
WHERE r.loan_history_id = @loan_history_id
ORDER BY r.renewal_date DESC;

-- name: GetRenewalsByDate :many
SELECT
    r.*,
    lh.book_id,
    b.title,
    b.author,
    rd.full_name as reader_name,
    lib.full_name as librarian_name
FROM renewals r
JOIN loan_history lh ON r.loan_history_id = lh.id
JOIN books b ON lh.book_id = b.id
JOIN readers rd ON lh.reader_id = rd.id
JOIN librarians lib ON r.librarian_id = lib.id
WHERE r.renewal_date BETWEEN @start_date AND @end_date
ORDER BY r.renewal_date DESC;

-- name: GetMostRenewedBooks :many
SELECT
    b.title,
    b.author,
    b.book_code,
    COUNT(r.id) as renewal_count
FROM renewals r
JOIN loan_history lh ON r.loan_history_id = lh.id
JOIN books b ON lh.book_id = b.id
WHERE r.renewal_date >= CURRENT_DATE - INTERVAL '@days_back days'
GROUP BY b.id, b.title, b.author, b.book_code
ORDER BY renewal_count DESC
LIMIT @result_limit;
