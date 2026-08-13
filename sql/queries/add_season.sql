-- name: AddSeason :one
INSERT INTO seasons(season_number, total_episodes, unwatched_episodes, finished, show_title)
VALUES (
    sqlc.arg(season_number),
    sqlc.arg(total_episodes),
    sqlc.arg(total_episodes),
    FALSE,
    sqlc.arg(show_title)
)
RETURNING *;