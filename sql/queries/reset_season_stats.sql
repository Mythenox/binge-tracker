-- name: ResetSeasonStats :exec
UPDATE episodes
SET viewtime = 0.0,
    watched = FALSE
WHERE 
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number);

UPDATE seasons
SET finished = FALSE,
    unwatched_episodes = total_episodes
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number);