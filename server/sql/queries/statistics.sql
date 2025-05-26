-- name: GetLoanStatusStatistics :many
SELECT
    status,
    COUNT(*) as count,
    ROUND(
        COUNT(*) * 100.0 / (
            SELECT
                COUNT(*)
            FROM
                loan_history lh
            WHERE
                lh.loan_date >= $1
        ),
        2
    ) as percentage
FROM
    loan_history
WHERE
    loan_date >= $1
GROUP BY
    status
ORDER BY
    count DESC;

-- name: CreateDailyStatistics :exec
INSERT INTO
    daily_statistics (
        stat_date,
        total_loans,
        total_returns,
        total_renewals,
        total_reservations,
        total_new_readers,
        total_fines_amount,
        overdue_books
    )
VALUES
    (
        CURRENT_DATE,
        (
            SELECT
                COUNT(*)
            FROM
                loan_history
            WHERE
                loan_date = CURRENT_DATE
        ),
        (
            SELECT
                COUNT(*)
            FROM
                loan_history
            WHERE
                return_date = CURRENT_DATE
        ),
        (
            SELECT
                COUNT(*)
            FROM
                renewals
            WHERE
                renewal_date = CURRENT_DATE
        ),
        (
            SELECT
                COUNT(*)
            FROM
                reservations
            WHERE
                reservation_date = CURRENT_DATE
        ),
        (
            SELECT
                COUNT(*)
            FROM
                readers
            WHERE
                registration_date = CURRENT_DATE
        ),
        (
            SELECT
                COALESCE(SUM(amount), 0)
            FROM
                fines
            WHERE
                fine_date = CURRENT_DATE
        ),
        (
            SELECT
                COUNT(*)
            FROM
                loan_history
            WHERE
                status = 'active'
                AND due_date < CURRENT_DATE
        )
    ) ON CONFLICT (stat_date) DO
UPDATE
SET
    total_loans = EXCLUDED.total_loans,
    total_returns = EXCLUDED.total_returns,
    total_renewals = EXCLUDED.total_renewals,
    total_reservations = EXCLUDED.total_reservations,
    total_new_readers = EXCLUDED.total_new_readers,
    total_fines_amount = EXCLUDED.total_fines_amount,
    overdue_books = EXCLUDED.overdue_books,
    updated_at = NOW ();

-- name: GetMonthlyReport :many
SELECT
    TO_CHAR (stat_date, 'YYYY-MM') as month,
    SUM(total_loans) as monthly_loans,
    SUM(total_returns) as monthly_returns,
    SUM(total_renewals) as monthly_renewals,
    SUM(total_new_readers) as new_readers,
    SUM(total_fines_amount) as fines_collected,
    AVG(overdue_books) as avg_overdue_books
FROM
    daily_statistics
WHERE
    stat_date >= CURRENT_DATE - INTERVAL '12 months'
GROUP BY
    TO_CHAR (stat_date, 'YYYY-MM')
ORDER BY
    month DESC;

-- name: GetYearlyReportByCategory :many
SELECT
    c.name as category,
    COUNT(DISTINCT b.id) as books_in_category,
    COUNT(lh.id) as total_loans_year,
    ROUND(
        AVG(
            EXTRACT(
                DAYS
                FROM
                    (lh.return_date - lh.loan_date)
            )
        ),
        1
    ) as avg_reading_days
FROM
    book_categories c
    LEFT JOIN books b ON c.id = b.category_id
    LEFT JOIN loan_history lh ON b.id = lh.book_id
    AND EXTRACT(
        YEAR
        FROM
            lh.loan_date
    ) = EXTRACT(
        YEAR
        FROM
            CURRENT_DATE
    )
GROUP BY
    c.id,
    c.name
ORDER BY
    total_loans_year DESC;

-- name: GetInventoryReport :many
SELECT
    h.name as hall,
    c.name as category,
    COUNT(b.id) as book_count,
    SUM(b.total_copies) as total_copies,
    SUM(b.available_copies) as available_copies,
    SUM(
        CASE
            WHEN b.status = 'lost' THEN 1
            ELSE 0
        END
    ) as lost_books,
    SUM(
        CASE
            WHEN b.condition_status = 'poor' THEN 1
            ELSE 0
        END
    ) as books_need_replacement
FROM
    halls h
    LEFT JOIN books b ON h.id = b.hall_id
    LEFT JOIN book_categories c ON b.category_id = c.id
GROUP BY
    h.id,
    h.name,
    c.id,
    c.name
ORDER BY
    h.name,
    c.name;
