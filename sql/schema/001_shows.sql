-- +goose Up
CREATE TABLE shows(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT UNIQUE NOT NULL,
    seasons INTEGER NOT NULL,
    total_episodes INTEGER NOT NULL,
    unwatched_episodes INTEGER NOT NULL
);

-- +goose Down
DROP TABLE shows;