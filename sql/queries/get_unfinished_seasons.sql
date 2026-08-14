-- name: GetUnfinishedSeasons :many
SELECT * FROM seasons
WHERE 
    finished = FALSE AND
    show_title = sqlc.arg(show_title)
ORDER BY season_number ASC;