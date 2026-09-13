-- name: CreateTrip :one
INSERT INTO trips (
    user_id,
    passport_id,
    destination_code,
    departure_date,
    return_date
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;


-- name: GetTripByID :one
SELECT *
FROM trips
WHERE id = $1;


-- name: ListTripsByUserID :many
SELECT *
FROM trips
WHERE user_id = $1
ORDER BY departure_date ASC, created_at DESC;


-- name: ListTripsByPassportID :many
SELECT *
FROM trips
WHERE passport_id = $1
ORDER BY departure_date ASC, created_at DESC;


-- name: UpdateTrip :one
UPDATE trips
SET
    passport_id = $2,
    destination_code = $3,
    departure_date = $4,
    return_date = $5,
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: DeleteTrip :exec
DELETE FROM trips
WHERE id = $1;