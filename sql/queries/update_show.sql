-- name: UpdateShow :one
UPDATE shows
SET
    unwatched_episodes = sqlc.arg(unwatched_episodes),
    total_episodes = sqlc.arg(total_episodes),
    seasons = sqlc.arg(seasons)
WHERE
    title = sqlc.arg(show_title)
RETURNING *;