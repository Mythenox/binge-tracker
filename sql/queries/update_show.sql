-- name: UpdateShow :one
UPDATE shows
SET
    unwatched_episodes = sqlc.arg(unwatched_episodes)
WHERE
    title = sqlc.arg(show_title)
RETURNING *;