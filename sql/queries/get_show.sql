-- name: GetShow :one
SELECT * FROM shows
WHERE title = sqlc.arg(show_title);