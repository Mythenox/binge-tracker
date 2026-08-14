-- name: GetSeasonsFromShow :many
SELECT * FROM seasons
WHERE show_title = sqlc.arg(show_title)
ORDER BY season_number ASC;