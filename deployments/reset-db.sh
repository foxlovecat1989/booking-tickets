#!/bin/bash

# Reset database to initial state
# This script runs every time the PostgreSQL container starts

echo "Resetting database to initial state..."

# Wait for PostgreSQL to be ready
until pg_isready -U postgres -d tickets_db; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

# Drop and recreate the database to ensure clean state
psql -U postgres -c "DROP DATABASE IF EXISTS tickets_db;"
psql -U postgres -c "CREATE DATABASE tickets_db;"

# Create the schema and sample data
psql -U postgres -d tickets_db << 'EOF'
-- Create concerts table
CREATE TABLE IF NOT EXISTS concerts (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  location VARCHAR(255) NOT NULL,
  description TEXT,
  created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000
);

-- Create concert_sessions table
CREATE TABLE IF NOT EXISTS concert_sessions (
  id SERIAL PRIMARY KEY,
  concert_id INTEGER NOT NULL,
  start_time BIGINT NOT NULL,
  end_time BIGINT NOT NULL,
  venue VARCHAR(255) NOT NULL,
  number_of_seats INTEGER NOT NULL DEFAULT 100,
  price DECIMAL(10,2) NOT NULL,
  FOREIGN KEY (concert_id) REFERENCES concerts(id) ON DELETE CASCADE
);

-- Create tickets table
CREATE TABLE IF NOT EXISTS tickets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id INTEGER NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('pending', 'sold', 'available')),
  FOREIGN KEY (session_id) REFERENCES concert_sessions(id) ON DELETE CASCADE
);

-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
  id SERIAL PRIMARY KEY,
  created_at BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW()) * 1000,
  status VARCHAR(50) NOT NULL DEFAULT 'pending',
  total_price DECIMAL(10,2) NOT NULL
);

-- Create schema_migrations table for migration tracking
CREATE TABLE IF NOT EXISTS schema_migrations (
  version BIGINT PRIMARY KEY,
  dirty BOOLEAN NOT NULL DEFAULT FALSE,
  applied_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_concert_sessions_concert_id ON concert_sessions(concert_id);
CREATE INDEX IF NOT EXISTS idx_tickets_session_id ON tickets(session_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);

-- Insert sample data
INSERT INTO concerts (name, location, description) VALUES 
('Rock Concert 2024', 'Madison Square Garden', 'The biggest rock concert of the year')
ON CONFLICT DO NOTHING;

-- Insert sample concert session (Dec 31, 2024)
INSERT INTO concert_sessions (concert_id, start_time, end_time, venue, number_of_seats, price) VALUES 
(1, 1735689600000, 1735693200000, 'Main Arena', 10000, 99.99)
ON CONFLICT DO NOTHING;

-- Insert sample tickets (10,000 available tickets)
INSERT INTO tickets (session_id, status)
SELECT 1, 'available'
FROM generate_series(1, 10000)
ON CONFLICT DO NOTHING;
EOF

echo "Database reset completed successfully!" 