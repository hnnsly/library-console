-- name: CreateOperationLog :exec
INSERT INTO operation_logs (
    operation_type, entity_type, entity_id, librarian_id, details, description
) VALUES (
    @operation_type, @entity_type, @entity_id, @librarian_id, @details, @description
);

-- name: GetOperationLogs :many
SELECT ol.*, l.full_name as librarian_name
FROM operation_logs ol
LEFT JOIN librarians l ON ol.librarian_id = l.id
WHERE (@entity_type::text IS NULL OR ol.entity_type = @entity_type)
  AND (@entity_id::int IS NULL OR ol.entity_id = @entity_id)
  AND ol.operation_date >= CURRENT_DATE - INTERVAL '@days_back days'
ORDER BY ol.operation_date DESC
LIMIT @page_limit OFFSET @page_offset;

-- name: GetRecentOperations :many
SELECT ol.*, l.full_name as librarian_name
FROM operation_logs ol
LEFT JOIN librarians l ON ol.librarian_id = l.id
ORDER BY ol.operation_date DESC
LIMIT @result_limit;
