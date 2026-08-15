-- name: ResetShowStats :exec
UPDATE shows
SET unwatched_episodes = total_episodes
WHERE title = sqlc.arg(show_title);

-- name: ResetSeasonStatsForShow :exec
UPDATE seasons
SET finished = FALSE,
    unwatched_episodes = total_episodes
WHERE show_title = sqlc.arg(show_title);

-- name: ResetEpisodeStatsForShow :exec
UPDATE episodes
SET viewtime = 0.0,
    watched = FALSE
WHERE show_title = sqlc.arg(show_title);