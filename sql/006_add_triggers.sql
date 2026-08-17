-- +goose Up
-- +goose StatementBegin
CREATE TRIGGER increment_on_episode_insertion
AFTER INSERT ON episodes
BEGIN
    UPDATE seasons
    SET
        unwatched_episodes = unwatched_episodes + 1,
        total_episodes = total_episodes + 1,
        finished = FALSE
    WHERE
        show_title = NEW.show_title AND
        season_number = NEW.season_number;
    
    UPDATE shows
    SET total_episodes = total_episodes + 1,
        unwatched_episodes = unwatched_episodes + 1
    WHERE title = NEW.show_title;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER increment_on_season_insertion
AFTER INSERT ON seasons
BEGIN
    UPDATE shows
    SET
        total_episodes = total_episodes + NEW.total_episodes,
        unwatched_episodes = unwatched_episodes + NEW.unwatched_episodes,
        seasons = seasons + 1
    WHERE title = NEW.show_title;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER increment_on_episode_insertion;
DROP TRIGGER increment_on_season_insertion;
-- +goose StatementEnd
    