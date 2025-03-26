-- name: AddBook :one
INSERT INTO books (title, author_id)
VALUES ($1, $2) RETURNING *;