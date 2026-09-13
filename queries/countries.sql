-- name: GetCountryByCode :one
SELECT *
FROM countries
WHERE code = $1;


-- name: ListCountries :many
SELECT *
FROM countries
ORDER BY name ASC;