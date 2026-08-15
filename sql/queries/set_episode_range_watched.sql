-- name: SetEpisodeRangeWatched :exec
UPDATE episodes
SET watched = TRUE,
    viewtime = runtime
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number) AND
    sqlc.arg(range_start) <= episode_number AND
    episode_number <= sqlc.arg(range_end);