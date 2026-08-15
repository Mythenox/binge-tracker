-- name: GetNextNoProgressEpisodeFromSeason :one
SELECT * FROM episodes
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number) AND
    viewtime = 0.0
ORDER BY episode_number ASC
LIMIT 1;