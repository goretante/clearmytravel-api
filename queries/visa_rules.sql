-- name: CreateVisaRule :one
INSERT INTO visa_rules (
    nationality_code,
    destination_code,
    visa_required,
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

-- name: GetVisaRuleByID :one
SELECT *
FROM visa_rules
WHERE id = $1;

-- name: ListVisaRules :many
SELECT *
FROM visa_rules
ORDER BY nationality_code, destination_code, effective_from;

-- name: FindVisaRule :one
SELECT *
FROM visa_rules
WHERE nationality_code = $1
  AND destination_code = $2
  AND effective_from <= $3
  AND (
    effective_to IS NULL
    OR effective_to >= $3
  )
ORDER BY effective_from DESC
LIMIT 1;

-- name: DeleteVisaRule :exec
DELETE FROM visa_rules
WHERE id = $1;
