-- name: CreateBook :one
INSERT INTO books (title, isbn, publication_year, publisher, pages, language, description, total_copies, available_copies)
VALUES (@title, @isbn, @publication_year, @publisher, @pages, @language, @description, @total_copies, @available_copies)
RETURNING *;

-- name: GetBookByID :one
SELECT * FROM books WHERE id = @book_id;

-- name: GetBookByISBN :one
SELECT * FROM books WHERE isbn = @isbn;

-- name: UpdateBook :one
UPDATE books
SET
    title = COALESCE(@title, title),
    isbn = COALESCE(@isbn, isbn),
    publication_year = COALESCE(@publication_year, publication_year),
    publisher = COALESCE(@publisher, publisher),
    pages = COALESCE(@pages, pages),
    language = COALESCE(@language, language),
    description = COALESCE(@description, description),
    total_copies = COALESCE(@total_copies, total_copies),
    available_copies = COALESCE(@available_copies, available_copies),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @book_id
RETURNING *;

-- name: DeleteBook :exec
DELETE FROM books WHERE id = @book_id;

-- name: ListBooks :many
SELECT * FROM books
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: SearchBooksByTitle :many
SELECT * FROM books
WHERE to_tsvector('russian', title) @@ plainto_tsquery('russian', @search_query)
ORDER BY title;

-- name: GetBooksWithAuthors :many
SELECT
    b.*,
    string_agg(a.full_name, ', ' ORDER BY a.full_name) as authors
FROM books b
LEFT JOIN book_authors ba ON b.id = ba.book_id
LEFT JOIN authors a ON ba.author_id = a.id
GROUP BY b.id, b.title, b.isbn, b.publication_year, b.publisher, b.pages, b.language, b.description, b.total_copies, b.available_copies, b.created_at, b.updated_at
ORDER BY b.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetBookWithDetails :one
SELECT
    b.*,
    string_agg(a.full_name, ', ' ORDER BY a.full_name) as authors,
    COALESCE(ROUND(AVG(br.rating), 2), 0) as avg_rating,
    COUNT(br.rating) as rating_count
FROM books b
LEFT JOIN book_authors ba ON b.id = ba.book_id
LEFT JOIN authors a ON ba.author_id = a.id
LEFT JOIN book_ratings br ON b.id = br.book_id
WHERE b.id = @book_id
GROUP BY b.id;

-- name: GetBooksByAuthor :many
SELECT b.*
FROM books b
JOIN book_authors ba ON b.id = ba.book_id
WHERE ba.author_id = @author_id
ORDER BY b.title;

-- name: UpdateBookAvailability :exec
UPDATE books
SET available_copies = @available_copies, updated_at = CURRENT_TIMESTAMP
WHERE id = @book_id;

-- name: CountBooks :one
SELECT COUNT(*) FROM books;

-- name: GetTopRatedBooks :many
SELECT
    b.*,
    string_agg(a.full_name, ', ' ORDER BY a.full_name) as authors,
    ROUND(AVG(br.rating), 2) as avg_rating,
    COUNT(br.rating) as rating_count
FROM books b
LEFT JOIN book_authors ba ON b.id = ba.book_id
LEFT JOIN authors a ON ba.author_id = a.id
JOIN book_ratings br ON b.id = br.book_id
GROUP BY b.id
HAVING COUNT(br.rating) >= @min_ratings
ORDER BY AVG(br.rating) DESC, COUNT(br.rating) DESC
LIMIT @limit_val;
