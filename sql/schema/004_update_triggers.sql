-- +goose Up

-- +goose StatementBegin
CREATE TRIGGER decrement_unwatched_episodes
AFTER UPDATE OF watched ON episodes
WHEN OLD.watched = 0 AND NEW.watched = 1 
BEGIN
    UPDATE seasons
    SET 
        unwatched_episodes = unwatched_episodes - 1,
        finished = IIF(unwatched_episodes - 1 <= 0, 1, 0)
    WHERE
        show_title = NEW.show_title AND
        season_number = NEW.season_number;
    
    UPDATE shows
    SET unwatched_episodes = unwatched_episodes - 1
    WHERE title = NEW.show_title;
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER increment_unwatched_episodes
AFTER UPDATE OF watched ON episodes
WHEN OLD.watched = 1 AND NEW.watched = 0
BEGIN
    UPDATE seasons
    SET 
        unwatched_episodes = unwatched_episodes + 1,
        finished = IIF(unwatched_episodes + 1 > 0, 0, 1)
    WHERE
        show_title = NEW.show_title AND
        season_number = NEW.season_number;
    
    UPDATE shows
    SET unwatched_episodes = unwatched_episodes + 1
    WHERE title = NEW.show_title;
END;
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS decrement_unwatched_episodes;
DROP TRIGGER IF EXISTS increment_unwatched_episodes;