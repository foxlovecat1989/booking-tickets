package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

// TestDBConfig holds database configuration for tests
type TestDBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// GetTestDBConfig returns database configuration for tests
func GetTestDBConfig() TestDBConfig {
	return TestDBConfig{
		Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		Port:     getEnvOrDefault("TEST_DB_PORT", "5432"),
		User:     getEnvOrDefault("TEST_DB_USER", "postgres"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "password"),
		DBName:   getEnvOrDefault("TEST_DB_NAME", "tickets_db"),
	}
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// SetupTestDB creates a test database connection and initializes schema
func SetupTestDB(t *testing.T) (*BaseRepository, func()) {
	config := GetTestDBConfig()

	// First connect to default postgres database to create test database if needed
	defaultDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		config.Host, config.Port, config.User, config.Password)

	defaultDB, err := sqlx.Connect("postgres", defaultDSN)
	if err != nil {
		t.Logf("Warning: Could not connect to default database: %v", err)
	} else {
		defer defaultDB.Close()

		// Create test database if it doesn't exist
		_, err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", config.DBName))
		if err != nil {
			// Database might already exist, which is fine
			t.Logf("Database creation result: %v", err)
		}
	}

	// Connect to test database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)

	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err)

	// Initialize schema
	err = initializeTestSchema(db)
	require.NoError(t, err)

	baseRepo := NewBaseRepository(db)

	// Clean up function
	cleanup := func() {
		db.Close()
	}

	return baseRepo, cleanup
}

// initializeTestSchema creates the necessary tables for testing
func initializeTestSchema(db *sqlx.DB) error {
	// Use a transaction to ensure atomic schema creation
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback() // Ignore rollback errors in defer
	}()

	// First, let's see what tables exist
	rows, err := tx.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
	`)
	if err != nil {
		return fmt.Errorf("failed to query existing tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	// Check if concert_sessions table exists and what columns it has
	concertSessionExists := false
	for _, table := range tables {
		if table == "concert_sessions" {
			concertSessionExists = true
			break
		}
	}

	if concertSessionExists {
		// Check if the table has the correct columns
		concertSessionRows, err := tx.Query(`
			SELECT column_name, data_type 
			FROM information_schema.columns 
			WHERE table_name = 'concert_sessions'
		`)
		if err == nil {
			defer concertSessionRows.Close()
			var hasNumberOfSeats bool
			for concertSessionRows.Next() {
				var columnName, dataType string
				if err := concertSessionRows.Scan(&columnName, &dataType); err != nil {
					continue
				}
				if columnName == "number_of_seats" {
					hasNumberOfSeats = true
					break
				}
			}
			if !hasNumberOfSeats {
				// Table exists but doesn't have the right columns, drop it
				_, err = tx.Exec("DROP TABLE IF EXISTS concert_sessions CASCADE")
				if err != nil {
					return fmt.Errorf("failed to drop concert_sessions table: %w", err)
				}
				concertSessionExists = false
			}
		}
	}

	// Only create schema if tables don't exist or are incomplete
	if !concertSessionExists {
		schema := `
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

		-- Note: order_items table removed as it's not used in the current application

		-- Create schema_migrations table for migration tests
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT NOW()
		);
		`

		// Execute schema creation
		_, err = tx.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Verify that the concert_sessions table was created with the correct columns
	var columnName string
	err = db.QueryRow(`
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_name = 'concert_sessions' AND column_name = 'number_of_seats'
	`).Scan(&columnName)

	if err != nil {
		return fmt.Errorf("number_of_seats column not found in concert_sessions table: %w", err)
	}

	return nil
}

// CleanupTestData cleans up test data from the database
func CleanupTestData(t *testing.T, baseRepo *BaseRepository) {
	// Clean up test data
	queries := []string{
		"DELETE FROM orders",
		"DELETE FROM tickets",
		"DELETE FROM concert_sessions",
		"DELETE FROM concerts",
		"DELETE FROM schema_migrations",
	}

	for _, query := range queries {
		_, err := baseRepo.db.Exec(query)
		if err != nil {
			t.Logf("Warning: Failed to cleanup test data: %v", err)
		}
	}
}
