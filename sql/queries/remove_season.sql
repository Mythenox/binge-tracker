-- name: RemoveSeason :exec
DELETE FROM seasons
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number);