-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER decrement_on_episode_deletion
AFTER DELETE ON episodes
BEGIN
    UPDATE seasons
    SET
        unwatched_episodes = IIF(OLD.watched = TRUE, unwatched_episodes - 1, unwatched_episodes),
        total_episodes = total_episodes - 1,
        finished = IIF(unwatched_episodes = total_episodes, TRUE, FALSE)
    WHERE
        show_title = OLD.show_title AND
        season_number = OLD.season_number;
    
    UPDATE shows
    SET total_episodes = total_episodes - 1,
        unwatched_episodes = IIF(OLD.watched = TRUE, unwatched_episodes - 1, unwatched_episodes)
    WHERE title = OLD.show_title;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER decrement_on_season_deletion
AFTER DELETE ON seasons
BEGIN
    UPDATE shows
    SET
        total_episodes = total_episodes - OLD.total_episodes,
        unwatched_episodes = unwatched_episodes - OLD.unwatched_episodes,
        seasons = seasons - 1
    WHERE title = OLD.show_title;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER decrement_on_episode_deletion;
DROP TRIGGER decrement_on_season_deletion;
-- +goose StatementEnd
    