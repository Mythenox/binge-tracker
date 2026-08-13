-- name: UpdateSeason :one
UPDATE seasons
SET
    unwatched_episodes = sqlc.arg(unwatched_episodes),
    finished = sqlc.arg(finished)
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number)
RETURNING *;