package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"tickets/internal/config"
	"tickets/internal/logger"
	"tickets/internal/migrations"

	_ "github.com/lib/pq"
)

func main() {
	var (
		command = flag.String("command", "up", "Migration command: up, down, status, create")
		steps   = flag.Int("steps", 1, "Number of migrations to rollback (for down command)")
		name    = flag.String("name", "", "Migration name (for create command)")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	if err := logger.Init(&cfg.Logging); err != nil {
		logger.Fatalf("Failed to initialize logger: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		logger.Fatalf("Failed to ping database: %v", err)
	}

	// Create migration manager
	manager := migrations.NewMigrationManager(db)

	// Load migrations
	if err := manager.LoadMigrations("migrations"); err != nil {
		logger.Fatalf("Failed to load migrations: %v", err)
	}

	switch *command {
	case "up":
		if err := manager.MigrateUp(); err != nil {
			logger.Fatalf("Failed to migrate up: %v", err)
		}
		logger.Info("Migrations applied successfully")

	case "down":
		if err := manager.MigrateDown(*steps); err != nil {
			logger.Fatalf("Failed to migrate down: %v", err)
		}
		logger.Infof("Rolled back %d migrations", *steps)

	case "status":
		status, err := manager.GetMigrationStatus()
		if err != nil {
			logger.Fatalf("Failed to get migration status: %v", err)
		}

		fmt.Printf("%-10s %-30s %-10s %-20s\n", "Version", "Name", "Applied", "Created At")
		fmt.Println(string(make([]byte, 80)))
		for _, s := range status {
			applied := "No"
			if s.Applied {
				applied = "Yes"
			}
			fmt.Printf("%-10d %-30s %-10s %-20s\n", s.Version, s.Name, applied, s.CreatedAt.Format("2006-01-02 15:04:05"))
		}

	case "create":
		if *name == "" {
			logger.Fatal("Migration name is required for create command")
		}
		if err := createMigration(*name); err != nil {
			logger.Fatalf("Failed to create migration: %v", err)
		}
		logger.Infof("Created migration: %s", *name)

	default:
		fmt.Println("Usage: migrate [options]")
		fmt.Println("Commands:")
		fmt.Println("  up     - Apply all pending migrations")
		fmt.Println("  down   - Rollback last N migrations")
		fmt.Println("  status - Show migration status")
		fmt.Println("  create - Create a new migration")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func createMigration(name string) error {
	// Find the next version number
	files, err := os.ReadDir("migrations")
	if err != nil {
		return err
	}

	maxVersion := 0
	for _, file := range files {
		if file.IsDir() || file.Name()[0] < '0' || file.Name()[0] > '9' {
			continue
		}

		var version int
		if _, err := fmt.Sscanf(file.Name(), "%d_", &version); err == nil {
			if version > maxVersion {
				maxVersion = version
			}
		}
	}

	nextVersion := maxVersion + 1

	// Create up migration file
	upFileName := fmt.Sprintf("migrations/%03d_%s.up.sql", nextVersion, name)
	upContent := fmt.Sprintf(`-- Migration: %s
-- Version: %d
-- Created: %s

-- Add your migration SQL here
-- Example:
-- CREATE TABLE example (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT NOW()
-- );

`, name, nextVersion, "now")

	if err := os.WriteFile(upFileName, []byte(upContent), 0644); err != nil {
		return err
	}

	// Create down migration file
	downFileName := fmt.Sprintf("migrations/%03d_%s.down.sql", nextVersion, name)
	downContent := fmt.Sprintf(`-- Rollback: %s
-- Version: %d
-- Created: %s

-- Add your rollback SQL here
-- Example:
-- DROP TABLE IF EXISTS example;

`, name, nextVersion, "now")

	return os.WriteFile(downFileName, []byte(downContent), 0644)
}
