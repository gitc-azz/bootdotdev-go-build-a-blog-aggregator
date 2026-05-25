-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetFeedsWithUserName :many
SELECT
    feeds.name,
    feeds.url,
    users.name as username
FROM feeds
INNER JOIN feed_follows ON feeds.id = feed_follows.feed_id
INNER JOIN users ON feed_follows.user_id = users.id;

-- name: GetFeed :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: ResetFeed :exec
DELETE FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
