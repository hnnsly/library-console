-- name: CreateAuthor :one
INSERT INTO authors (full_name, birth_year, death_year, biography)
VALUES (@full_name, @birth_year, @death_year, @biography)
RETURNING *;

-- name: GetAuthorByID :one
SELECT * FROM authors WHERE id = @author_id;

-- name: UpdateAuthor :one
UPDATE authors
SET
    full_name = COALESCE(@full_name, full_name),
    birth_year = COALESCE(@birth_year, birth_year),
    death_year = COALESCE(@death_year, death_year),
    biography = COALESCE(@biography, biography)
WHERE id = @author_id
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = @author_id;

-- name: ListAuthors :many
SELECT * FROM authors ORDER BY full_name
LIMIT @limit_val OFFSET @offset_val;

-- name: SearchAuthorsByName :many
SELECT * FROM authors
WHERE to_tsvector('russian', full_name) @@ plainto_tsquery('russian', @search_query)
ORDER BY full_name;

-- name: GetAuthorsByBook :many
SELECT a.*
FROM authors a
JOIN book_authors ba ON a.id = ba.author_id
WHERE ba.book_id = @book_id
ORDER BY a.full_name;

-- name: CountAuthors :one
SELECT COUNT(*) FROM authors;
