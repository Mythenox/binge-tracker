-- name: RemoveEpisode :exec
DELETE FROM episodes
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number) AND
    episode_number = sqlc.arg(episode_number);