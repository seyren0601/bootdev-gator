-- name: CreatePost :one
INSERT INTO posts(created_at, updated_at, title, url, description, published_at, feed_id)
VALUES(
    $1,
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT
    feeds.name AS source,
    posts.title,
    posts.url,
    posts.description,
    posts.published_at
FROM posts
    INNER JOIN feeds ON feeds.id = posts.feed_id
    INNER JOIN feed_follows ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2;