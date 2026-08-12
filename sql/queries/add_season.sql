-- name: AddSeason :one
INSERT INTO seasons(season_number, total_episodes, unwatched_episodes, finished, show_title)
VALUES (
    sql.arg(season_number),
    sqlc.arg(total_episodes),
    sqlc.arg(total_episodes),
    FALSE,
    sqlc.arg(show_title)
)
RETURNING *;