-- name: CreateBookCopy :one
INSERT INTO book_copies (book_id, copy_code, status, reading_hall_id, condition_notes)
VALUES (@book_id, @copy_code, @status, @reading_hall_id, @condition_notes)
RETURNING *;

-- name: GetBookCopyByID :one
SELECT bc.*, b.title, b.isbn
FROM book_copies bc
JOIN books b ON bc.book_id = b.id
WHERE bc.id = @copy_id;

-- name: GetBookCopyByCode :one
SELECT bc.*, b.title, b.isbn
FROM book_copies bc
JOIN books b ON bc.book_id = b.id
WHERE bc.copy_code = @copy_code;

-- name: UpdateBookCopy :one
UPDATE book_copies
SET
    status = COALESCE(@status, status),
    reading_hall_id = COALESCE(@reading_hall_id, reading_hall_id),
    condition_notes = COALESCE(@condition_notes, condition_notes),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @copy_id
RETURNING *;

-- name: DeleteBookCopy :exec
DELETE FROM book_copies WHERE id = @copy_id;

-- name: ListBookCopiesByBook :many
SELECT bc.*, rh.hall_name
FROM book_copies bc
LEFT JOIN reading_halls rh ON bc.reading_hall_id = rh.id
WHERE bc.book_id = @book_id
ORDER BY bc.copy_code;

-- name: ListAvailableBookCopies :many
SELECT bc.*, b.title, b.isbn
FROM book_copies bc
JOIN books b ON bc.book_id = b.id
WHERE bc.status = 'available'
ORDER BY b.title, bc.copy_code;

-- name: GetBookCopiesByStatus :many
SELECT bc.*, b.title, b.isbn, rh.hall_name
FROM book_copies bc
JOIN books b ON bc.book_id = b.id
LEFT JOIN reading_halls rh ON bc.reading_hall_id = rh.id
WHERE bc.status = @status
ORDER BY b.title, bc.copy_code;

-- name: CountBookCopiesByBook :one
SELECT COUNT(*) FROM book_copies WHERE book_id = @book_id;

-- name: CountAvailableBookCopies :one
SELECT COUNT(*) FROM book_copies WHERE book_id = @book_id AND status = 'available';
