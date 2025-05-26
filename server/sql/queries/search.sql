-- name: GlobalSearch :many
(SELECT
    'book' as type,
    b.id,
    b.title as name,
    b.author as details,
    NULL as ticket
FROM books b
WHERE b.title ILIKE '%' || @search_term::text || '%'
   OR b.author ILIKE '%' || @search_term::text || '%')
UNION ALL
(SELECT
    'reader' as type,
    r.id,
    r.full_name as name,
    r.phone as details,
    r.ticket_number as ticket
FROM readers r
WHERE r.full_name ILIKE '%' || @search_term::text || '%'
   OR r.ticket_number ILIKE '%' || @search_term::text || '%')
ORDER BY name
LIMIT 20;

-- name: AdvancedBookSearch :many
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
