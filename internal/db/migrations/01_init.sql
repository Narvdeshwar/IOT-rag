CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE devices (
    id TEXT PRIMARY KEY,
    name TEXT,
    zone TEXT,
    type TEXT,
    last_seen TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE sensor_events (
    id BIGSERIAL PRIMARY KEY,
    device_id TEXT REFERENCES devices(id),
    event_time TIMESTAMPTZ DEFAULT NOW(),
    metric TEXT,
    value DOUBLE PRECISION,
    unit TEXT,
    metadata JSONB
);

CREATE TABLE telemetry_chunks (
    id BIGSERIAL PRIMARY KEY,
    content TEXT,
    device_id TEXT,
    event_time TIMESTAMPTZ,
    metadata JSONB,
    embedding VECTOR(1536)
);

CREATE INDEX ON telemetry_chunks USING hnsw (embedding vector_cosine_ops);
