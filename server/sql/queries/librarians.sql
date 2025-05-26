-- name: CreateLibrarian :one
INSERT INTO librarians (
    full_name, employee_id, position, phone, email, hire_date
) VALUES (
    @full_name, @employee_id, @position, @phone, @email, CURRENT_DATE
) RETURNING *;

-- name: GetLibrarianByID :one
SELECT * FROM librarians WHERE id = @librarian_id;

-- name: GetLibrarianByEmployeeID :one
SELECT * FROM librarians WHERE employee_id = @employee_id;

-- name: GetAllLibrarians :many
SELECT * FROM librarians WHERE status = 'active' ORDER BY full_name;

-- name: UpdateLibrarian :one
UPDATE librarians
SET full_name = @full_name,
    position = @position,
    phone = @phone,
    email = @email,
    updated_at = NOW()
WHERE id = @librarian_id
RETURNING *;

-- name: DeactivateLibrarian :exec
UPDATE librarians
SET status = 'inactive', updated_at = NOW()
WHERE id = @librarian_id;
