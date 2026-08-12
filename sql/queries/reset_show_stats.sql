-- name: ResetShowStats :exec
UPDATE episodes
SET viewtime = 0.0,
    watched = FALSE
WHERE show_title = sqlc.arg(show_title);

UPDATE seasons
SET finished = FALSE,
    unwatched_episodes = total_episodes
WHERE show_title = sqlc.arg(show_title);

UPDATE shows
SET unwatched_episodes = total_episodes
WHERE title = sqlc.arg(show_title);