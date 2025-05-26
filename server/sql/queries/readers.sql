-- name: CreateReader :one
INSERT INTO readers (
    full_name, ticket_number, birth_date, phone, email,
    education, hall_id, registration_date
) VALUES (
    @full_name, @ticket_number, @birth_date, @phone, @email,
    @education, @hall_id, CURRENT_DATE
) RETURNING *;

-- name: GetReaderByTicket :one
SELECT r.*, h.name as hall_name, h.specialization
FROM readers r
JOIN halls h ON r.hall_id = h.id
WHERE r.ticket_number = @ticket_number;

-- name: GetReaderByID :one
SELECT r.*, h.name as hall_name, h.specialization
FROM readers r
JOIN halls h ON r.hall_id = h.id
WHERE r.id = @reader_id;

-- name: SearchReadersByName :many
SELECT r.*, h.name as hall_name
FROM readers r
JOIN halls h ON r.hall_id = h.id
WHERE r.full_name ILIKE '%' || @search_name || '%'
ORDER BY r.full_name
LIMIT @page_limit OFFSET @page_offset;

-- name: UpdateReader :one
UPDATE readers
SET full_name = @full_name,
    phone = @phone,
    email = @email,
    education = @education,
    hall_id = @hall_id,
    updated_at = NOW()
WHERE id = @reader_id
RETURNING *;

-- name: UpdateReaderStatus :exec
UPDATE readers
SET status = @status, updated_at = NOW()
WHERE id = @reader_id;

-- name: GetAllReaders :many
SELECT r.*, h.name as hall_name, h.specialization,
       (SELECT COUNT(*) FROM loan_history lh WHERE lh.reader_id = r.id AND lh.status = 'active') as current_loans
FROM readers r
JOIN halls h ON r.hall_id = h.id
ORDER BY r.full_name
LIMIT @page_limit OFFSET @page_offset;

-- name: GetReadersCount :one
SELECT COUNT(*) FROM readers;

-- name: UpdateReaderDebt :exec
UPDATE readers r
SET total_debt = (
    SELECT COALESCE(SUM(f.amount), 0)
    FROM fines f
    WHERE f.reader_id = readers.id AND f.status = 'unpaid'
),
updated_at = NOW()
WHERE r.id = @reader_id;

-- name: GetReaderStatistics :one
SELECT
    r.full_name,
    r.ticket_number,
    COUNT(lh.id) as total_books_taken,
    COUNT(CASE WHEN lh.status = 'active' THEN 1 END) as current_loans,
    COUNT(CASE WHEN lh.return_date > lh.due_date THEN 1 END) as late_returns,
    ROUND(AVG(CASE WHEN lh.return_date IS NOT NULL
              THEN EXTRACT(DAYS FROM (lh.return_date - lh.loan_date)) END), 1) as avg_reading_days,
    COALESCE(SUM(f.amount), 0) as total_fines,
    COALESCE(SUM(CASE WHEN f.status = 'unpaid' THEN f.amount ELSE 0 END), 0) as unpaid_fines,
    MAX(lh.loan_date) as last_activity
FROM readers r
LEFT JOIN loan_history lh ON r.id = lh.reader_id
LEFT JOIN fines f ON r.id = f.reader_id
WHERE r.id = @reader_id
GROUP BY r.id, r.full_name, r.ticket_number;

-- name: GetReaderFavoriteCategories :many
SELECT c.name as category_name, COUNT(*) as books_count,
       ROUND(COUNT(*) * 100.0 / (SELECT COUNT(*) FROM loan_history WHERE lh.reader_id = @reader_id), 2) as percentage
FROM loan_history lh
JOIN books b ON lh.book_id = b.id
JOIN book_categories c ON b.category_id = c.id
WHERE lh.reader_id = @reader_id
GROUP BY c.id, c.name
ORDER BY books_count DESC
LIMIT 5;

-- name: GetActiveReaders :many
SELECT
    r.full_name,
    r.ticket_number,
    COUNT(lh.id) as total_loans,
    COUNT(CASE WHEN lh.status = 'active' THEN 1 END) as current_loans,
    MAX(lh.loan_date) as last_loan_date
FROM readers r
LEFT JOIN loan_history lh ON r.id = lh.reader_id
WHERE lh.loan_date >= CURRENT_DATE - INTERVAL '@days_back days' OR lh.loan_date IS NULL
GROUP BY r.id, r.full_name, r.ticket_number
HAVING COUNT(lh.id) > 0
ORDER BY total_loans DESC
LIMIT @result_limit;
