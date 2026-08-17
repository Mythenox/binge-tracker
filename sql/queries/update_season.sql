-- name: UpdateSeason :one
UPDATE seasons
SET
    unwatched_episodes = sqlc.arg(unwatched_episodes),
    finished = sqlc.arg(finished),
    total_episodes = sqlc.arg(total_episodes)
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number)
RETURNING *;