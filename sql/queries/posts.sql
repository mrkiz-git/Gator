-- Select all posts
-- name: GetAllPosts :many
SELECT * FROM posts 
ORDER BY published_at DESC
LIMIT $1;

-- Select a single post by ID
-- name: GetPostByID :one
SELECT * FROM posts WHERE id = $1;

-- Insert a new post with default timestamps for created_at and updated_at
-- name: CreatePost :one
INSERT INTO posts (created_at, updated_at, id,title, url, description, published_at, feed_id)
VALUES (NOW(), NOW(), $1, $2, $3, $4, $5, $6)
RETURNING *;

-- Update an existing post
-- name: UpdatePost :one
UPDATE posts
SET title = $1, url = $2, description = $3, published_at = $4, feed_id = $5, updated_at = now()
WHERE id = $6
RETURNING *;

-- Delete a post by ID
-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;
