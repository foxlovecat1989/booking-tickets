-- Rollback: initial_schema
-- Version: 1
-- Created: 2024-01-01

-- Drop all tables in reverse order to handle foreign key constraints

-- Drop indexes first (only the ones that were created)
DROP INDEX IF EXISTS idx_tickets_status;
DROP INDEX IF EXISTS idx_tickets_session_id;
DROP INDEX IF EXISTS idx_concert_sessions_concert_id;

-- Drop tables
DROP TABLE IF EXISTS orders CASCADE;
DROP TABLE IF EXISTS tickets CASCADE;
DROP TABLE IF EXISTS concert_sessions CASCADE;
DROP TABLE IF EXISTS concerts CASCADE; 