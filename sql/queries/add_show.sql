-- name: AddShow :one
INSERT INTO shows(title, seasons, total_episodes, unwatched_episodes)
VALUES (
    sqlc.arg(title),
    1,
    sqlc.arg(total_episodes),
    sqlc.arg(total_episodes)
)
RETURNING *;