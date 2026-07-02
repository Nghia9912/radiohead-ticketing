-- Enable extension to generate UUIDs
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    total_capacity INT NOT NULL,
    status VARCHAR(50) DEFAULT 'UPCOMING'
);

CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID REFERENCES events(id),
    seat_identifier VARCHAR(50) NOT NULL,
    status VARCHAR(50) DEFAULT 'AVAILABLE',
    version INT DEFAULT 1,
    UNIQUE(event_id, seat_identifier) -- Strict constraint to prevent duplicate seats
);

-- Optimize querying speed for available tickets
CREATE INDEX idx_tickets_availability ON tickets(event_id, status) WHERE status = 'AVAILABLE';