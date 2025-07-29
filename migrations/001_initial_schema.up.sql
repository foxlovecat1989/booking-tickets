-- Migration: initial_schema
-- Version: 1
-- Created: 2024-01-01

-- Initialize the database with concert booking schema
-- Note: Database is already created by Docker environment variables

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

-- Note: order_items and payments tables removed as they are not used in the current application

-- Create indexes for better performance
-- Only create indexes that are actually used by queries
CREATE INDEX IF NOT EXISTS idx_concert_sessions_concert_id ON concert_sessions(concert_id);
CREATE INDEX IF NOT EXISTS idx_tickets_session_id ON tickets(session_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);