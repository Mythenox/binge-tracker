-- name: ResetEpisodeStats :exec
UPDATE episodes
SET viewtime = 0.0,
    watched = FALSE
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number) AND
    episode_number = sqlc.arg(episode_number);