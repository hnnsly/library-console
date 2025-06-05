-- name: CreateBookRating :one
INSERT INTO book_ratings (book_id, reader_id, rating, review, rating_date)
VALUES (@book_id, @reader_id, @rating, @review, @rating_date)
ON CONFLICT (book_id, reader_id)
DO UPDATE SET
    rating = EXCLUDED.rating,
    review = EXCLUDED.review,
    rating_date = EXCLUDED.rating_date
RETURNING *;

-- name: GetBookRatingByID :one
SELECT
    br.*,
    r.full_name as reader_name,
    b.title as book_title
FROM book_ratings br
JOIN readers r ON br.reader_id = r.id
JOIN books b ON br.book_id = b.id
WHERE br.id = @rating_id;

-- name: GetReaderBookRating :one
SELECT * FROM book_ratings
WHERE book_id = @book_id AND reader_id = @reader_id;

-- name: UpdateBookRating :one
UPDATE book_ratings
SET
    rating = COALESCE(@rating, rating),
    review = COALESCE(@review, review),
    rating_date = CURRENT_DATE
WHERE id = @rating_id
RETURNING *;

-- name: DeleteBookRating :exec
DELETE FROM book_ratings WHERE id = @rating_id;

-- name: GetBookRatings :many
SELECT
    br.*,
    r.full_name as reader_name
FROM book_ratings br
JOIN readers r ON br.reader_id = r.id
WHERE br.book_id = @book_id
ORDER BY br.rating_date DESC;

-- name: GetReaderRatings :many
SELECT
    br.*,
    b.title as book_title
FROM book_ratings br
JOIN books b ON br.book_id = b.id
WHERE br.reader_id = @reader_id
ORDER BY br.rating_date DESC;

-- name: GetBookAverageRating :one
SELECT
    COALESCE(ROUND(AVG(rating), 2), 0) as avg_rating,
    COUNT(*) as total_ratings
FROM book_ratings
WHERE book_id = @book_id;

-- name: GetTopRatedBooksWithRatings :many
SELECT
    b.id,
    b.title,
    string_agg(a.full_name, ', ' ORDER BY a.full_name) as authors,
    ROUND(AVG(br.rating), 2) as avg_rating,
    COUNT(br.rating) as rating_count
FROM books b
LEFT JOIN book_authors ba ON b.id = ba.book_id
LEFT JOIN authors a ON ba.author_id = a.id
JOIN book_ratings br ON b.id = br.book_id
GROUP BY b.id, b.title
HAVING COUNT(br.rating) >= @min_ratings
ORDER BY AVG(br.rating) DESC, COUNT(br.rating) DESC
LIMIT @limit_val;
