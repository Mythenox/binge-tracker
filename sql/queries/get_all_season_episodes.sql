-- name: GetAllSeasonEpisodes :many
SELECT * FROM episodes
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number)
ORDER BY episode_number ASC;