-- name: CreateAuthor :one
INSERT INTO authors (full_name)
VALUES (@full_name)
RETURNING id, full_name;

-- name: GetOrCreateAuthor :one
INSERT INTO authors (full_name)
VALUES (@full_name)
ON CONFLICT (full_name) DO UPDATE SET full_name = EXCLUDED.full_name
RETURNING id, full_name;

-- name: GetAuthorById :one
SELECT id, full_name
FROM authors
WHERE id = @id;

-- name: SearchAuthors :many
SELECT id, full_name
FROM authors
WHERE full_name ILIKE '%' || @search_term || '%'
ORDER BY full_name;

-- name: GetAllAuthors :many
SELECT id, full_name
FROM authors
ORDER BY full_name;
