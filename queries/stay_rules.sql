-- name: CreateStayRule :one
INSERT INTO stay_rules (
    nationality_code,
    destination_code,
    max_stay_days,
    effective_from,
    effective_to,
    source_name,
    source_url
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetStayRuleByID :one
SELECT *
FROM stay_rules
WHERE id = $1;

-- name: ListStayRules :many
SELECT *
FROM stay_rules
ORDER BY nationality_code, destination_code, effective_from;

-- name: FindStayRule :one
SELECT *
FROM stay_rules
WHERE nationality_code = $1
  AND destination_code = $2
  AND effective_from <= $3
  AND (
      effective_to IS NULL
      OR effective_to >= $3
  )
ORDER BY effective_from DESC
LIMIT 1;

-- name: DeleteStayRule :exec
DELETE FROM stay_rules
WHERE id = $1;