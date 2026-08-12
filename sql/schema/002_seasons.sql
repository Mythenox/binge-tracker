-- +goose Up
CREATE TABLE seasons(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    season_number INTEGER NOT NULL,
    total_episodes INTEGER NOT NULL,
    unwatched_episodes INTEGER NOT NULL,
    finished BOOLEAN NOT NULL,
    show_id UUID NOT NULL,
    FOREIGN KEY (show_id)
    REFERENCES shows (id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE seasons;