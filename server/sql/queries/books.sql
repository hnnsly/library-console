-- name: CreateBook :one
INSERT INTO books (
    title, author, publication_year, isbn, book_code,
    category_id, hall_id, total_copies, available_copies,
    condition_status, location_info, acquisition_date
) VALUES (
    @title, @author, @publication_year, @isbn, @book_code,
    @category_id, @hall_id, @total_copies, @available_copies,
    @condition_status, @location_info, CURRENT_DATE
) RETURNING *;

-- name: GetBookByID :one
SELECT b.*, c.name as category_name, h.name as hall_name
FROM books b
LEFT JOIN book_categories c ON b.category_id = c.id
JOIN halls h ON b.hall_id = h.id
WHERE b.id = @book_id;

-- name: GetBookByCode :one
SELECT b.*, c.name as category_name, h.name as hall_name
FROM books b
LEFT JOIN book_categories c ON b.category_id = c.id
JOIN halls h ON b.hall_id = h.id
WHERE b.book_code = @book_code;

-- name: SearchBooks :many
SELECT b.*, c.name as category_name, h.name as hall_name
FROM books b
LEFT JOIN book_categories c ON b.category_id = c.id
JOIN halls h ON b.hall_id = h.id
WHERE
    (@title::text IS NULL OR b.title ILIKE '%' || @title || '%')
    AND (@author::text IS NULL OR b.author ILIKE '%' || @author || '%')
    AND (@book_code::text IS NULL OR b.book_code = @book_code)
    AND (@isbn::text IS NULL OR b.isbn = @isbn)
    AND (@category_id::int IS NULL OR b.category_id = @category_id)
    AND (@hall_id::int IS NULL OR b.hall_id = @hall_id)
    AND b.status != 'lost'
ORDER BY b.title
LIMIT @page_limit OFFSET @page_offset;

-- name: GetAvailableBooks :many
SELECT b.*, c.name as category_name, h.name as hall_name
FROM books b
LEFT JOIN book_categories c ON b.category_id = c.id
JOIN halls h ON b.hall_id = h.id
WHERE b.available_copies > 0 AND b.status = 'available'
ORDER BY b.popularity_score DESC
LIMIT @result_limit;

-- name: UpdateBookCopies :exec
UPDATE books
SET total_copies = @total_copies,
    available_copies = @available_copies,
    updated_at = NOW()
WHERE id = @book_id;

-- name: WriteOffBook :exec
UPDATE books
SET total_copies = total_copies - 1,
    available_copies = GREATEST(0, available_copies - 1),
    updated_at = NOW()
WHERE id = @book_id AND total_copies > 0;

-- name: UpdateBookAvailability :exec
UPDATE books
SET available_copies = @available_copies,
    popularity_score = @popularity_score,
    updated_at = NOW()
WHERE id = @book_id;

-- name: GetBooksByAuthorInHall :one
SELECT COUNT(*) as book_count,
       SUM(b.total_copies) as total_copies,
       SUM(b.available_copies) as available_copies
FROM books b
WHERE b.author ILIKE '%' || @author || '%' AND b.hall_id = @hall_id;

-- name: GetBooksWithSingleCopy :many
SELECT DISTINCT r.full_name, r.ticket_number, r.phone, b.title, b.author
FROM readers r
JOIN loan_history lh ON r.id = lh.reader_id
JOIN books b ON lh.book_id = b.id
WHERE lh.status = 'active' AND b.total_copies = 1
ORDER BY r.full_name;

-- name: GetTopRatedBooks :many
SELECT title, author, rating, popularity_score, book_code
FROM books
WHERE rating = (SELECT MAX(rating) FROM books)
ORDER BY popularity_score DESC
LIMIT @result_limit;

-- name: GetPopularBooks :many
SELECT
    b.title,
    b.author,
    b.book_code,
    COUNT(lh.id) as loan_count,
    b.popularity_score,
    AVG(b.rating) as avg_rating
FROM books b
JOIN loan_history lh ON b.id = lh.book_id
WHERE lh.loan_date >= CURRENT_DATE - INTERVAL '@days_back days'
GROUP BY b.id, b.title, b.author, b.book_code, b.popularity_score
ORDER BY loan_count DESC
LIMIT @result_limit;

-- name: AdvancedSearchBooks :many
SELECT
    b.id,
    b.title,
    b.author,
    b.publication_year,
    b.book_code,
    b.isbn,
    c.name as category,
    h.name as hall,
    b.total_copies,
    b.available_copies,
    b.popularity_score,
    b.rating,
    CASE WHEN b.available_copies > 0 THEN 'Доступна' ELSE 'Недоступна' END as availability_status
FROM books b
LEFT JOIN book_categories c ON b.category_id = c.id
JOIN halls h ON b.hall_id = h.id
WHERE
    (@title_filter::text = '' OR b.title ILIKE '%' || @title_filter || '%')
    AND (@author_filter::text = '' OR b.author ILIKE '%' || @author_filter || '%')
    AND (@year_filter::int = 0 OR b.publication_year = @year_filter)
    AND (@category_filter::int = 0 OR b.category_id = @category_filter)
    AND (@hall_filter::int = 0 OR b.hall_id = @hall_filter)
    AND (@available_only::boolean = false OR b.available_copies > 0)
    AND b.status != 'lost'
ORDER BY
    CASE @sort_by::text
        WHEN 'title' THEN b.title
        WHEN 'author' THEN b.author
        WHEN 'year' THEN b.publication_year::text
        WHEN 'popularity' THEN b.popularity_score::text
        ELSE b.title
    END
LIMIT @page_limit OFFSET @page_offset;
