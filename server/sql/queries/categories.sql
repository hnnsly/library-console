
-- name: GetAllCategories :many
SELECT * FROM book_categories ORDER BY name;

-- name: CreateCategory :one
INSERT INTO book_categories (name, description, default_loan_days)
VALUES (@name, @description, @default_loan_days)
RETURNING *;

-- name: GetCategoryStatistics :many
SELECT
    c.name as category_name,
    COUNT(DISTINCT b.id) as total_books,
    COUNT(lh.id) as total_loans,
    ROUND(AVG(CASE WHEN lh.return_date IS NOT NULL
              THEN EXTRACT(DAYS FROM (lh.return_date - lh.loan_date)) END), 1) as avg_loan_days,
    COUNT(CASE WHEN lh.return_date > lh.due_date THEN 1 END) as overdue_returns,
    ROUND(COUNT(CASE WHEN lh.return_date > lh.due_date THEN 1 END) * 100.0 /
          NULLIF(COUNT(CASE WHEN lh.return_date IS NOT NULL THEN 1 END), 0), 2) as overdue_percentage
FROM book_categories c
LEFT JOIN books b ON c.id = b.category_id
LEFT JOIN loan_history lh ON b.id = lh.book_id AND lh.loan_date >= CURRENT_DATE - INTERVAL '@days_back days'
GROUP BY c.id, c.name
ORDER BY total_loans DESC;
