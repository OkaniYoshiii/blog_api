-- name: ListPosts :many
SELECT id, title, content FROM posts;

-- name: CreatePost :exec
INSERT INTO posts (
   title, content
) VALUES (
    ?, ?
);
