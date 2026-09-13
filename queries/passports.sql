-- name: CreatePassport :one
INSERT INTO passports (
    user_id,
    nationality_code,
    expires_at
)
VALUES ($1, $2, $3)
RETURNING *;


-- name: GetPassportByID :one
SELECT *
FROM passports
WHERE id = $1;


-- name: ListPassportsByUserID :many
SELECT *
FROM passports
WHERE user_id = $1
ORDER BY created_at DESC;


-- name: UpdatePassport :one
UPDATE passports
SET
    nationality_code = $2,
    expires_at = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: DeletePassport :exec
DELETE FROM passports
WHERE id = $1;