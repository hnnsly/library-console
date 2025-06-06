-- name: IssueBook :one
INSERT INTO book_issues (reader_id, book_copy_id, due_date, librarian_id)
VALUES (@reader_id, @book_copy_id, @due_date, @librarian_id)
RETURNING id, issue_date, due_date;

-- name: ReturnBook :one
UPDATE book_issues
SET return_date = CURRENT_DATE
WHERE book_copy_id = @book_copy_id AND return_date IS NULL
RETURNING id, reader_id, book_copy_id, issue_date, due_date, return_date;

-- name: GetBooksToReturn :many
SELECT
    bi.id,
    r.ticket_number,
    r.full_name as reader_name,
    b.title,
    bc.copy_code,
    bi.issue_date,
    bi.due_date,
    CASE
        WHEN bi.due_date < CURRENT_DATE THEN (CURRENT_DATE - bi.due_date)
        ELSE 0
    END as days_overdue
FROM book_issues bi
JOIN readers r ON bi.reader_id = r.id
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
WHERE bi.return_date IS NULL
ORDER BY bi.due_date;

-- name: GetOverdueBooks :many
SELECT
    bi.id,
    r.ticket_number,
    r.full_name as reader_name,
    b.title,
    bc.copy_code,
    bi.issue_date,
    bi.due_date,
    (CURRENT_DATE - bi.due_date) as days_overdue
FROM book_issues bi
JOIN readers r ON bi.reader_id = r.id
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
WHERE bi.return_date IS NULL AND bi.due_date < CURRENT_DATE
ORDER BY bi.due_date;

-- name: GetReaderActiveBooks :many
SELECT
    bi.id,
    b.title,
    bc.copy_code,
    bi.issue_date,
    bi.due_date
FROM book_issues bi
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
WHERE bi.reader_id = @reader_id AND bi.return_date IS NULL
ORDER BY bi.due_date;

-- name: GetRecentBookOperations :many
SELECT
    'issue' as operation_type,
    bi.created_at as operation_time,
    r.full_name as reader_name,
    r.ticket_number,
    b.title as book_title,
    bc.copy_code,
    u.username as librarian_name,
    bi.due_date::text as additional_info
FROM book_issues bi
JOIN readers r ON bi.reader_id = r.id
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
LEFT JOIN users u ON bi.librarian_id = u.id
WHERE bi.created_at >= @since_date

UNION ALL

SELECT
    'return' as operation_type,
    bi.updated_at as operation_time,
    r.full_name as reader_name,
    r.ticket_number,
    b.title as book_title,
    bc.copy_code,
    u.username as librarian_name,
    bi.return_date::text as additional_info
FROM book_issues bi
JOIN readers r ON bi.reader_id = r.id
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
LEFT JOIN users u ON bi.librarian_id = u.id
WHERE bi.return_date IS NOT NULL AND bi.updated_at >= @since_date

ORDER BY operation_time DESC
LIMIT @limit_count;
