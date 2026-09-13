-- name: CreateETARule :one
INSERT INTO eta_rules (
    nationality_code,
    destination_code,
    eta_required,
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

-- name: GetETARuleByID :one
SELECT *
FROM eta_rules
WHERE id = $1;

-- name: ListETARules :many
SELECT *
FROM eta_rules
ORDER BY nationality_code, destination_code, effective_from;

-- name: FindETARule :one
SELECT *
FROM eta_rules
WHERE nationality_code = $1
  AND destination_code = $2
  AND effective_from <= $3
  AND (
      effective_to IS NULL
      OR effective_to >= $3
  )
ORDER BY effective_from DESC
LIMIT 1;

-- name: DeleteETARule :exec
DELETE FROM eta_rules
WHERE id = $1;