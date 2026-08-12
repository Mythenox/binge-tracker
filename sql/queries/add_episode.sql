-- name: AddEpisode :one
INSERT INTO episodes(
    episode_number,
    filepath,
    viewtime,
    runtime,
    watched,
    show_title,
    season_number         
)
VALUES (
    sqlc.arg(episode_number),
    sqlc.arg(filepath),
    0.0,
    sqlc.arg(runtime),
    FALSE,
    sqlc.arg(show_title),
    sqlc.arg(season_number)
)
RETURNING *;