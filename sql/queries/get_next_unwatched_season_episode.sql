-- name: GetNextUnwatchedEpisodeFromSeason :one
SELECT *
FROM episodes
WHERE 
    watched = FALSE AND
    season_number = sqlc.arg(season_number) AND
    show_title = sqlc.arg(show_title)
ORDER BY episode_number ASC
LIMIT 1;