-- name: CreateFeed :one
INSERT INTO feeds (id, user_id, name, url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, name, url, created_at, updated_at;

-- name: GetFeedByID :many
SELECT id, user_id, name, url, created_at, updated_at
FROM feeds
WHERE id = $1;

-- name: ListFeeds :many
SELECT 
    feeds.id, 
    feeds.user_id,
    users.name user_name,
    feeds.name, 
    feeds.url, 
    feeds.created_at, 
    feeds.updated_at
FROM feeds 
INNER JOIN users 
ON feeds.user_id = users.id
ORDER BY feeds.created_at DESC;

-- name: ListFeedsByUserID :many
SELECT id, user_id, name, url, created_at, updated_at
FROM feeds
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateFeed :one
UPDATE feeds
SET name = $2, url = $3, updated_at = $4
WHERE id = $1
RETURNING id, user_id, name, url, created_at, updated_at;

-- name: DeleteFeed :exec
DELETE FROM feeds
WHERE id = $1;


-- name: GetFeedByURL :many
SELECT * FROM feeds
WHERE URL = $1;

-- name: MarkFeedFetched :one
UPDATE feeds s SET last_fetched_at = $1
WHERE id = $2
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;