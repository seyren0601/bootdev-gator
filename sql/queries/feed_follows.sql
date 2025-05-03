-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows(user_id, feed_id, created_at, updated_at)
    VALUES(
        $1,
        $2,
        $3,
        $4
    )
    RETURNING *
)
SELECT  inserted.user_id AS UserID,
        inserted.feed_id AS FeedID,
        users.name AS UserName,
        feeds.name AS FeedName
FROM inserted
    INNER JOIN feeds ON inserted.feed_id = feeds.id
    INNER JOIN users ON inserted.user_id = users.id;

-- name: GetFeedFollowsForUser :many
SELECT  users.name AS user_name,
        feeds.name AS feed_name
FROM feed_follows
    INNER JOIN users ON users.id = feed_follows.user_id
    INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollowForUser :exec
DELETE 
FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;