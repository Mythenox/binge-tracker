-- name: GetAllShowEpisodes :many
SELECT * FROM episodes
WHERE show_title = sqlc.arg(show_title);