CREATE TABLE IF NOT EXISTS pastes (
    id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    max_views INT,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    views INT NOT NULL DEFAULT 0
);
