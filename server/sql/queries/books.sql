-- name: CreateBook :one
INSERT INTO books (title, isbn, publication_year, publisher, total_copies, available_copies)
VALUES (@title, @isbn, @publication_year, @publisher, @total_copies, @total_copies)
RETURNING id, title, isbn, publication_year, publisher, total_copies, available_copies;

-- name: UpdateBook :one
UPDATE books
SET title = @title, isbn = @isbn, publication_year = @publication_year, publisher = @publisher
WHERE id = @id
RETURNING id, title, isbn, publication_year, publisher, total_copies, available_copies;

-- name: GetBookById :one
SELECT id, title, isbn, publication_year, publisher, total_copies, available_copies
FROM books
WHERE id = @id;

-- name: SearchBooks :many
SELECT DISTINCT
    b.id,
    b.title,
    b.isbn,
    b.publication_year,
    b.publisher,
    b.available_copies,
    b.total_copies,
    STRING_AGG(a.full_name, ', ') as authors
FROM books b
LEFT JOIN book_authors ba ON b.id = ba.book_id
LEFT JOIN authors a ON ba.author_id = a.id
WHERE
    (@title::text IS NULL OR b.title ILIKE '%' || @title || '%') AND
    (@author::text IS NULL OR a.full_name ILIKE '%' || @author || '%') AND
    (@publication_year::int IS NULL OR b.publication_year = @publication_year)
GROUP BY b.id, b.title, b.isbn, b.publication_year, b.publisher, b.available_copies, b.total_copies
ORDER BY b.title;

-- name: GetAllBooks :many
SELECT DISTINCT
    b.id,
    b.title,
    b.isbn,
    b.publication_year,
    b.publisher,
    b.available_copies,
    b.total_copies,
    STRING_AGG(a.full_name, ', ') as authors
FROM books b
LEFT JOIN book_authors ba ON b.id = ba.book_id
LEFT JOIN authors a ON ba.author_id = a.id
GROUP BY b.id, b.title, b.isbn, b.publication_year, b.publisher, b.available_copies, b.total_copies
ORDER BY b.title;

-- name: UpdateBookCopies :exec
UPDATE books
SET available_copies = available_copies + @change
WHERE id = @book_id;
