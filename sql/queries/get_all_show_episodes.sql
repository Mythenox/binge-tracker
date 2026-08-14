-- name: GetAllShowEpisodes :many
SELECT * FROM episodes
WHERE show_title = sqlc.arg(show_title)
ORDER BY season_number ASC, episode_number ASC;