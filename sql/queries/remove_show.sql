-- name: RemoveShow :exec
DELETE FROM shows
WHERE title = sqlc.arg(show_title);