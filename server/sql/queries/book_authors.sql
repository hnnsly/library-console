-- name: AddBookAuthor :exec
INSERT INTO book_authors (book_id, author_id)
VALUES (@book_id, @author_id);

-- name: RemoveBookAuthor :exec
DELETE FROM book_authors
WHERE book_id = @book_id AND author_id = @author_id;

-- name: GetBookAuthors :many
SELECT a.id, a.full_name
FROM authors a
JOIN book_authors ba ON a.id = ba.author_id
WHERE ba.book_id = @book_id
ORDER BY a.full_name;

-- name: GetAuthorBooks :many
SELECT b.id, b.title, b.isbn, b.publication_year
FROM books b
JOIN book_authors ba ON b.id = ba.book_id
WHERE ba.author_id = @author_id
ORDER BY b.title;
