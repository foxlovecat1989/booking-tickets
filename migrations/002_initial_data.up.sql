-- Migration: initial_data
-- Version: 2
-- Created: 2024-01-01

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