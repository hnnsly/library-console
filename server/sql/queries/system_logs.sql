-- name: CreateSystemLog :exec
INSERT INTO system_logs (user_id, action_type, entity_type, entity_id, old_values, new_values, ip_address, user_agent, details)
VALUES (@user_id, @action_type, @entity_type, @entity_id, @old_values, @new_values, @ip_address, @user_agent, @details);

-- name: GetSystemLogByID :one
SELECT
    sl.*,
    u.username
FROM system_logs sl
LEFT JOIN users u ON sl.user_id = u.id
WHERE sl.id = @log_id;

-- name: GetSystemLogsByUser :many
SELECT
    sl.*,
    u.username
FROM system_logs sl
LEFT JOIN users u ON sl.user_id = u.id
WHERE sl.user_id = @user_id
ORDER BY sl.action_timestamp DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetSystemLogsByAction :many
SELECT
    sl.*,
    u.username
FROM system_logs sl
LEFT JOIN users u ON sl.user_id = u.id
WHERE sl.action_type = @action_type
ORDER BY sl.action_timestamp DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetSystemLogsByEntity :many
SELECT
    sl.*,
    u.username
FROM system_logs sl
LEFT JOIN users u ON sl.user_id = u.id
WHERE sl.entity_type = @entity_type AND sl.entity_id = @entity_id
ORDER BY sl.action_timestamp DESC;

-- name: GetRecentSystemLogs :many
SELECT
    sl.*,
    u.username
FROM system_logs sl
LEFT JOIN users u ON sl.user_id = u.id
WHERE sl.action_timestamp >= @from_timestamp
ORDER BY sl.action_timestamp DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetSystemLogsInDateRange :many
SELECT
    sl.*,
    u.username
FROM system_logs sl
LEFT JOIN users u ON sl.user_id = u.id
WHERE sl.action_timestamp BETWEEN @from_timestamp AND @to_timestamp
ORDER BY sl.action_timestamp DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: CountSystemLogsByAction :one
SELECT COUNT(*) FROM system_logs
WHERE action_type = @action_type
AND action_timestamp BETWEEN @from_timestamp AND @to_timestamp;

-- name: GetSystemLogsStatistics :one
SELECT
    COUNT(*) as total_logs,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(CASE WHEN action_type = 'login' THEN 1 END) as login_count,
    COUNT(CASE WHEN action_type = 'book_issue' THEN 1 END) as book_issue_count,
    COUNT(CASE WHEN action_type = 'book_return' THEN 1 END) as book_return_count
FROM system_logs
WHERE action_timestamp BETWEEN @from_timestamp AND @to_timestamp;
