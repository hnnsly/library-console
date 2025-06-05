-- name: GetLibraryStatistics :one
SELECT
    (SELECT COUNT(*) FROM books) as total_books,
    (SELECT COUNT(*) FROM book_copies) as total_copies,
    (SELECT COUNT(*) FROM book_copies WHERE status = 'available') as available_copies,
    (SELECT COUNT(*) FROM readers WHERE is_active = true) as active_readers,
    (SELECT COUNT(*) FROM book_issues WHERE return_date IS NULL) as active_issues,
    (SELECT COUNT(*) FROM book_issues WHERE return_date IS NULL AND due_date < CURRENT_DATE) as overdue_issues,
    (SELECT COUNT(*) FROM fines WHERE is_paid = false) as unpaid_fines,
    (SELECT COALESCE(SUM(amount - paid_amount), 0) FROM fines WHERE is_paid = false) as total_debt;

-- name: GetPopularBooks :many
SELECT
    b.id,
    b.title,
    string_agg(a.full_name, ', ' ORDER BY a.full_name) as authors,
    COUNT(bi.id) as issue_count,
    COALESCE(ROUND(AVG(br.rating), 2), 0) as avg_rating
FROM books b
LEFT JOIN book_authors ba ON b.id = ba.book_id
LEFT JOIN authors a ON ba.author_id = a.id
LEFT JOIN book_copies bc ON b.id = bc.book_id
LEFT JOIN book_issues bi ON bc.id = bi.book_copy_id
LEFT JOIN book_ratings br ON b.id = br.book_id
GROUP BY b.id, b.title
ORDER BY COUNT(bi.id) DESC, AVG(br.rating) DESC NULLS LAST
LIMIT @limit_val;

-- name: GetActiveReaderStatistics :many
SELECT
    r.id,
    r.full_name,
    r.ticket_number,
    rh.hall_name,
    COUNT(bi.id) as active_books,
    COUNT(CASE WHEN bi.due_date < CURRENT_DATE THEN 1 END) as overdue_books,
    COALESCE(SUM(f.amount - f.paid_amount), 0) as total_debt
FROM readers r
LEFT JOIN reading_halls rh ON r.reading_hall_id = rh.id
LEFT JOIN book_issues bi ON r.id = bi.reader_id AND bi.return_date IS NULL
LEFT JOIN fines f ON r.id = f.reader_id AND f.is_paid = false
WHERE r.is_active = true
GROUP BY r.id, r.full_name, r.ticket_number, rh.hall_name
HAVING COUNT(bi.id) > 0 OR COALESCE(SUM(f.amount - f.paid_amount), 0) > 0
ORDER BY COUNT(bi.id) DESC, total_debt DESC;

-- name: GetBooksByAuthorInHall :many
SELECT
    b.id,
    b.title,
    bc.copy_code,
    bc.status,
    rh.hall_name
FROM books b
JOIN book_authors ba ON b.id = ba.book_id
JOIN authors a ON ba.author_id = a.id
JOIN book_copies bc ON b.id = bc.book_id
LEFT JOIN reading_halls rh ON bc.reading_hall_id = rh.id
WHERE a.id = @author_id
AND (@hall_id::uuid IS NULL OR bc.reading_hall_id = @hall_id)
ORDER BY b.title, bc.copy_code;

-- name: GetReadersWithSingleCopyBooks :many
SELECT DISTINCT
    r.id,
    r.full_name,
    r.ticket_number,
    b.title as book_title,
    bc.copy_code
FROM readers r
JOIN book_issues bi ON r.id = bi.reader_id
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
WHERE bi.return_date IS NULL
AND b.total_copies = 1
ORDER BY r.full_name, b.title;

-- name: GetMonthlyStatistics :many
SELECT
    DATE_TRUNC('month', bi.issue_date) as month,
    COUNT(*) as issues_count,
    COUNT(DISTINCT bi.reader_id) as unique_readers,
    COUNT(DISTINCT bc.book_id) as unique_books
FROM book_issues bi
JOIN book_copies bc ON bi.book_copy_id = bc.id
WHERE bi.issue_date >= @from_date AND bi.issue_date <= @to_date
GROUP BY DATE_TRUNC('month', bi.issue_date)
ORDER BY month;

-- name: GetHallUtilizationReport :many
SELECT
    rh.id,
    rh.hall_name,
    rh.total_seats,
    rh.occupied_seats,
    COUNT(bc.id) as books_in_hall,
    COUNT(CASE WHEN bc.status = 'available' THEN 1 END) as available_books,
    COUNT(r.id) as registered_readers,
    COUNT(CASE WHEN r.is_active THEN 1 END) as active_readers
FROM reading_halls rh
LEFT JOIN book_copies bc ON rh.id = bc.reading_hall_id
LEFT JOIN readers r ON rh.id = r.reading_hall_id
GROUP BY rh.id, rh.hall_name, rh.total_seats, rh.occupied_seats
ORDER BY rh.hall_name;
