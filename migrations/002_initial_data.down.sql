-- Migration: initial_data (DOWN)
-- Version: 2
-- Created: 2024-01-01

-- Remove all tickets for the concert session
DELETE FROM tickets WHERE session_id = 1;

-- Remove the concert session
DELETE FROM concert_sessions WHERE id = 1;

-- Remove the concert
DELETE FROM concerts WHERE id = 1; 