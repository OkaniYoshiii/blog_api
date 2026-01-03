-- name: ListPosts :many
SELECT id, title, content FROM posts;

-- name: CreatePost :one
INSERT INTO posts (
   title, content
) VALUES (
    ?, ?
) RETURNING *;
