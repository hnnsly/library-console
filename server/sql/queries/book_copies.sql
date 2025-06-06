-- name: CreateBookCopy :one
INSERT INTO book_copies (book_id, copy_code, hall_id, location_info)
VALUES (@book_id, @copy_code, @hall_id, @location_info)
RETURNING id, copy_code, status;

-- name: UpdateBookCopyStatus :exec
UPDATE book_copies
SET status = @status
WHERE id = @copy_id;

-- name: GetBookCopyById :one
SELECT bc.id, bc.copy_code, bc.status, b.title, rh.hall_name, bc.location_info
FROM book_copies bc
JOIN books b ON bc.book_id = b.id
LEFT JOIN reading_halls rh ON bc.hall_id = rh.id
WHERE bc.id = @copy_id;

-- name: GetBookCopyByCode :one
SELECT bc.id, bc.copy_code, bc.status, b.id as book_id, b.title, bc.location_info
FROM book_copies bc
JOIN books b ON bc.book_id = b.id
WHERE bc.copy_code = @copy_code;

-- name: GetAvailableBookCopy :one
SELECT bc.id, bc.copy_code, bc.status, b.title, bc.location_info
FROM book_copies bc
JOIN books b ON bc.book_id = b.id
WHERE bc.copy_code = @copy_code AND bc.status = 'available';

-- name: GetBookCopiesByBookId :many
SELECT id, copy_code, status
FROM book_copies
WHERE book_id = @book_id
ORDER BY copy_code;

-- name: GetBookCopiesByHall :many
SELECT bc.id, bc.copy_code, bc.status, b.title
FROM book_copies bc
JOIN books b ON bc.book_id = b.id
WHERE bc.hall_id = @hall_id
ORDER BY b.title, bc.copy_code;
