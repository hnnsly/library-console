-- name: CreateBookIssue :one
INSERT INTO book_issues (reader_id, book_copy_id, issue_date, due_date, librarian_id, notes)
VALUES (@reader_id, @book_copy_id, @issue_date, @due_date, @librarian_id, @notes)
RETURNING *;

-- name: GetBookIssueByID :one
SELECT
    bi.*,
    r.full_name as reader_name,
    r.ticket_number,
    b.title as book_title,
    bc.copy_code,
    u.username as librarian_name
FROM book_issues bi
JOIN readers r ON bi.reader_id = r.id
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
LEFT JOIN users u ON bi.librarian_id = u.id
WHERE bi.id = @issue_id;

-- name: UpdateBookIssue :one
UPDATE book_issues
SET
    due_date = COALESCE(@due_date, due_date),
    return_date = COALESCE(@return_date, return_date),
    extended_count = COALESCE(@extended_count, extended_count),
    notes = COALESCE(@notes, notes),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @issue_id
RETURNING *;

-- name: ReturnBook :exec
UPDATE book_issues
SET return_date = @return_date, updated_at = CURRENT_TIMESTAMP
WHERE id = @issue_id;

-- name: ExtendBookIssue :exec
UPDATE book_issues
SET
    due_date = @new_due_date,
    extended_count = extended_count + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = @issue_id;

-- name: GetActiveIssuesByReader :many
SELECT
    bi.*,
    b.title as book_title,
    bc.copy_code,
    CASE
        WHEN bi.due_date < CURRENT_DATE THEN (CURRENT_DATE - bi.due_date)
        ELSE 0
    END as overdue_days
FROM book_issues bi
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
WHERE bi.reader_id = @reader_id AND bi.return_date IS NULL
ORDER BY bi.due_date;

-- name: GetAllActiveIssues :many
SELECT
    bi.*,
    r.full_name as reader_name,
    r.ticket_number,
    b.title as book_title,
    bc.copy_code,
    CASE
        WHEN bi.due_date < CURRENT_DATE THEN (CURRENT_DATE - bi.due_date)
        ELSE 0
    END as overdue_days
FROM book_issues bi
JOIN readers r ON bi.reader_id = r.id
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
WHERE bi.return_date IS NULL
ORDER BY bi.due_date;

-- name: GetOverdueIssues :many
SELECT
    bi.*,
    r.full_name as reader_name,
    r.ticket_number,
    b.title as book_title,
    bc.copy_code,
    (CURRENT_DATE - bi.due_date) as overdue_days
FROM book_issues bi
JOIN readers r ON bi.reader_id = r.id
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
WHERE bi.return_date IS NULL AND bi.due_date < CURRENT_DATE
ORDER BY bi.due_date;

-- name: GetIssueHistory :many
SELECT
    bi.*,
    b.title as book_title,
    bc.copy_code
FROM book_issues bi
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
WHERE bi.reader_id = @reader_id
ORDER BY bi.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetBookIssueHistory :many
SELECT
    bi.*,
    r.full_name as reader_name,
    r.ticket_number
FROM book_issues bi
JOIN readers r ON bi.reader_id = r.id
WHERE bi.book_copy_id = @copy_id
ORDER BY bi.created_at DESC;

-- name: CountActiveIssuesByReader :one
SELECT COUNT(*) FROM book_issues
WHERE reader_id = @reader_id AND return_date IS NULL;

-- name: CountOverdueIssues :one
SELECT COUNT(*) FROM book_issues
WHERE return_date IS NULL AND due_date < CURRENT_DATE;

-- name: GetIssueDueSoon :many
SELECT
    bi.*,
    r.full_name as reader_name,
    r.ticket_number,
    b.title as book_title,
    bc.copy_code
FROM book_issues bi
JOIN readers r ON bi.reader_id = r.id
JOIN book_copies bc ON bi.book_copy_id = bc.id
JOIN books b ON bc.book_id = b.id
WHERE bi.return_date IS NULL
AND bi.due_date BETWEEN CURRENT_DATE AND (CURRENT_DATE + INTERVAL '@days_ahead days')
ORDER BY bi.due_date;
