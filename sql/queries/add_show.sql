-- name: AddShow :one
INSERT INTO shows(title, seasons, total_episodes, unwatched_episodes)
VALUES (
    sqlc.arg(title),
    sqlc.arg(seasons),
    sqlc.arg(total_episodes),
    sqlc.arg(total_episodes)
)
RETURNING *;