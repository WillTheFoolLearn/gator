-- name: FeedsUser :one
SELECT * FROM users
WHERE ID = $1;