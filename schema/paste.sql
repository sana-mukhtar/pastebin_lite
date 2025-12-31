CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS pastes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content TEXT NOT NULL,
    ttl_seconds INT,
    max_views INT,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    views INT DEFAULT 0
);
