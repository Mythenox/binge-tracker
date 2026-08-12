-- name: UpdateViewTime :one
UPDATE episodes
SET viewtime = sqlc.arg(viewtime),
    watched = sqlc.arg(watched)
WHERE
    show_title = sqlc.arg(show_title) AND
    season_number = sqlc.arg(season_number) AND
    episode_number = sqlc.arg(episode_number)
RETURNING *;