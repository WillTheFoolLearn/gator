-- name: FeedFromUrl :one
SELECT * FROM feeds
WHERE url = $1;