CREATE TABLE sequence (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    open_tracking_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    click_tracking_enabled BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE step (
    id SERIAL PRIMARY KEY,
    sequence_id INTEGER NOT NULL,
    subject VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (sequence_id) REFERENCES sequence (id)
);