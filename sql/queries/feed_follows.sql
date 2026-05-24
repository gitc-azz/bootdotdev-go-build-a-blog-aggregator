-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM inserted
INNER JOIN users ON users.id = inserted.user_id
INNER JOIN feeds ON feeds.id = inserted.feed_id;


-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM feed_follows
INNER JOIN users on users.id = feed_follows.user_id
INNER JOIN feeds on feeds.id = feed_follows.feed_id
WHERE users.name = $1;


-- name: ResetFeedFollows :exec
DELETE FROM feed_follows;

-- name: RmFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 and feed_id = $2;
