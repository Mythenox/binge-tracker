-- +goose Up
CREATE TABLE seasons(
    season_number INTEGER NOT NULL,
    total_episodes INTEGER NOT NULL,
    unwatched_episodes INTEGER NOT NULL,
    finished BOOLEAN NOT NULL,
    show_title TEXT NOT NULL,
    PRIMARY KEY (show_title, season_number),
    FOREIGN KEY (show_title)
    REFERENCES shows (title)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE seasons;