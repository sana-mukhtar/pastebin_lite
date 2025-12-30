CREATE TABLE pastes (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    ttl_seconds INT,
    max_views INT,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    views INT DEFAULT 0
);
