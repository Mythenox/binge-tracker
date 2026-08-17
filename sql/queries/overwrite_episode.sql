-- name: OverwriteEpisode :one
UPDATE episodes
SET viewtime = 0,
    watched = FALSE,
    runtime = sqlc.arg(runtime),
    filepath = sqlc.arg(filepath)
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number) AND
    episode_number = sqlc.arg(episode_number)
RETURNING *;