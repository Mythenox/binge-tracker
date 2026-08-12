-- +goose Up
CREATE TABLE episodes(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    episode_number INTEGER UNIQUE NOT NULL,
    viewtime FLOAT NOT NULL,
    runtime FLOAT NOT NULL,
    watched BOOLEAN NOT NULL,
    season_id UUID NOT NULL,
    FOREIGN KEY (season_id)
    REFERENCES seasons (id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE episodes;