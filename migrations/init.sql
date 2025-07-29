-- Initialize the database with concert booking schema
-- Note: Database is already created by Docker environment variables
\c tickets_db;

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

-- Create order_items table
CREATE TABLE IF NOT EXISTS order_items (
  id SERIAL PRIMARY KEY,
  order_id INTEGER NOT NULL,
  ticket_id UUID UNIQUE NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
  FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE CASCADE
);

-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
  id SERIAL PRIMARY KEY,
  order_id INTEGER NOT NULL,
  payment_method VARCHAR(100) NOT NULL,
  paid_at BIGINT,
  amount DECIMAL(10,2) NOT NULL,
  status VARCHAR(50) NOT NULL DEFAULT 'pending',
  transaction_id VARCHAR(255),
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_concerts_name ON concerts(name);
CREATE INDEX IF NOT EXISTS idx_concert_sessions_concert_id ON concert_sessions(concert_id);
CREATE INDEX IF NOT EXISTS idx_concert_sessions_start_time ON concert_sessions(start_time);
CREATE INDEX IF NOT EXISTS idx_tickets_session_id ON tickets(session_id);
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_payments_order_id ON payments(order_id);

-- Insert initial data
-- Insert initial concert data
INSERT INTO concerts (name, location, description) VALUES (
  'Rock Concert 2024',
  'Madison Square Garden, New York',
  'An amazing rock concert featuring top artists performing live music. Experience the energy and excitement of live rock music in one of the most iconic venues in the world.'
);

-- Insert concert session for the concert
INSERT INTO concert_sessions (concert_id, start_time, end_time, venue, number_of_seats, price) VALUES (
  1, -- concert_id from the inserted concert
  EXTRACT(EPOCH FROM '2024-12-31 20:00:00'::timestamp) * 1000, -- start_time: Dec 31, 2024 at 8 PM
  EXTRACT(EPOCH FROM '2024-12-31 23:00:00'::timestamp) * 1000, -- end_time: Dec 31, 2024 at 11 PM
  'Main Arena',
  10000, -- number_of_seats
  99.99 -- price per ticket
);

-- Insert 10,000 available tickets for the concert session
-- Using a more efficient approach with generate_series
INSERT INTO tickets (session_id, status)
SELECT 
  1, -- session_id from the inserted concert session
  'available'
FROM generate_series(1, 10000);

