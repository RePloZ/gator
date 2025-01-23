-- name: CreateFeeds :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: GetFeedsInformation :many
SELECT feeds.name AS feed_name,
    feeds.url,
    users.name AS user_name
FROM feeds
    JOIN users ON feeds.user_id = users.id;
-- name: GetFeeedByUrl :one
SELECT *
FROM feeds
WHERE feeds.url = $1;
-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $2,
    last_fetched_at = $2
WHERE id = $1;
-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;