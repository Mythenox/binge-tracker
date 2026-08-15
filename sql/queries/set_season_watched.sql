-- name: SetSeasonFinished :exec
UPDATE seasons
SET finished = TRUE
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number);

-- name: SetSeasonEpisodesWatched :exec

UPDATE episodes
SET
    watched = TRUE,
    viewtime = runtime
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number);