-- name: CreateUser :one
INSERT INTO users(created_at, updated_at, name)
VALUES(
    $1,
    $2,
    $3
)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE name = $1;

-- name: GetUserFromId :one
SELECT *
FROM users
WHERE id = $1;

-- name: DatabaseReset :exec
WITH CTE1 AS(
    DELETE FROM users
), CTE2 AS (
    DELETE FROM feeds
)
SELECT * FROM users;

-- name: GetUsers :many
SELECT *
FROM users;