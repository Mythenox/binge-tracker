-- +goose Up
CREATE TABLE episodes(
    episode_number INTEGER NOT NULL,
    filepath TEXT UNIQUE NOT NULL,
    viewtime FLOAT NOT NULL,
    runtime FLOAT NOT NULL,
    watched BOOLEAN NOT NULL,
    show_title TEXT NOT NULL,
    season_number INTEGER NOT NULL,
    PRIMARY KEY (show_title, season_number, episode_number),
    FOREIGN KEY (show_title, season_number)
    REFERENCES seasons (show_title, season_number)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE episodes;