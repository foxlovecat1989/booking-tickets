package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"tickets/internal/logger"

	_ "github.com/lib/pq"
)

// Migration represents a database migration
type Migration struct {
	Version   int64
	Name      string
	UpSQL     string
	DownSQL   string
	CreatedAt time.Time
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db         *sql.DB
	migrations []Migration
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// LoadMigrations loads migration files from the migrations directory
func (mm *MigrationManager) LoadMigrations(migrationsPath string) error {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Parse migration file name: 001_initial_schema.up.sql
		parts := strings.Split(strings.TrimSuffix(file.Name(), ".sql"), ".")
		if len(parts) != 2 {
			continue
		}

		versionName := parts[0]
		direction := parts[1]

		// Extract version number
		versionParts := strings.SplitN(versionName, "_", 2)
		if len(versionParts) != 2 {
			continue
		}

		version, err := strconv.ParseInt(versionParts[0], 10, 64)
		if err != nil {
			continue
		}

		name := versionParts[1]

		// Read migration content
		content, err := os.ReadFile(filepath.Join(migrationsPath, file.Name()))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		// Find or create migration
		var migration *Migration
		for i := range mm.migrations {
			if mm.migrations[i].Version == version && mm.migrations[i].Name == name {
				migration = &mm.migrations[i]
				break
			}
		}

		if migration == nil {
			mm.migrations = append(mm.migrations, Migration{
				Version:   version,
				Name:      name,
				CreatedAt: time.Now(),
			})
			migration = &mm.migrations[len(mm.migrations)-1]
		}

		// Set SQL content based on direction
		switch direction {
		case "up":
			migration.UpSQL = string(content)
		case "down":
			migration.DownSQL = string(content)
		}
	}

	// Sort migrations by version
	sort.Slice(mm.migrations, func(i, j int) bool {
		return mm.migrations[i].Version < mm.migrations[j].Version
	})

	return nil
}

// InitializeMigrationTable creates the migrations table if it doesn't exist
func (mm *MigrationManager) InitializeMigrationTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			dirty BOOLEAN NOT NULL DEFAULT FALSE,
			applied_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`
	_, err := mm.db.Exec(query)
	return err
}

// GetAppliedMigrations returns a list of applied migration versions
func (mm *MigrationManager) GetAppliedMigrations() (map[int64]bool, error) {
	applied := make(map[int64]bool)

	query := `SELECT version FROM schema_migrations ORDER BY version`
	rows, err := mm.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, nil
}

// ApplyMigration applies a single migration
func (mm *MigrationManager) ApplyMigration(migration Migration) error {
	// Check if migration is already applied
	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return err
	}

	if applied[migration.Version] {
		logger.Infof("Migration %d_%s already applied, skipping", migration.Version, migration.Name)
		return nil
	}

	// Begin transaction
	tx, err := mm.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Errorf("Failed to rollback transaction: %v", err)
		}
	}()

	// Execute migration
	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return fmt.Errorf("failed to apply migration %d_%s: %w", migration.Version, migration.Name, err)
	}

	// Record migration
	recordQuery := `INSERT INTO schema_migrations (version, applied_at) VALUES ($1, NOW())`
	if _, err := tx.Exec(recordQuery, migration.Version); err != nil {
		return fmt.Errorf("failed to record migration %d_%s: %w", migration.Version, migration.Name, err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	logger.Infof("Applied migration %d_%s", migration.Version, migration.Name)
	return nil
}

// RollbackMigration rolls back a single migration
func (mm *MigrationManager) RollbackMigration(migration Migration) error {
	// Check if migration is applied
	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return err
	}

	if !applied[migration.Version] {
		logger.Infof("Migration %d_%s not applied, skipping rollback", migration.Version, migration.Name)
		return nil
	}

	// Begin transaction
	tx, err := mm.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			logger.Errorf("Failed to rollback transaction: %v", err)
		}
	}()

	// Execute rollback
	if _, err := tx.Exec(migration.DownSQL); err != nil {
		return fmt.Errorf("failed to rollback migration %d_%s: %w", migration.Version, migration.Name, err)
	}

	// Remove migration record
	recordQuery := `DELETE FROM schema_migrations WHERE version = $1`
	if _, err := tx.Exec(recordQuery, migration.Version); err != nil {
		return fmt.Errorf("failed to remove migration record %d_%s: %w", migration.Version, migration.Name, err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	logger.Infof("Rolled back migration %d_%s", migration.Version, migration.Name)
	return nil
}

// MigrateUp applies all pending migrations
func (mm *MigrationManager) MigrateUp() error {
	if err := mm.InitializeMigrationTable(); err != nil {
		return fmt.Errorf("failed to initialize migration table: %w", err)
	}

	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return err
	}

	for _, migration := range mm.migrations {
		if !applied[migration.Version] {
			if err := mm.ApplyMigration(migration); err != nil {
				return err
			}
		}
	}

	logger.Info("All migrations applied successfully")
	return nil
}

// MigrateDown rolls back the last N migrations
func (mm *MigrationManager) MigrateDown(steps int) error {
	if err := mm.InitializeMigrationTable(); err != nil {
		return fmt.Errorf("failed to initialize migration table: %w", err)
	}

	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return err
	}

	// Find applied migrations in reverse order
	var appliedMigrations []Migration
	for _, migration := range mm.migrations {
		if applied[migration.Version] {
			appliedMigrations = append(appliedMigrations, migration)
		}
	}

	// Sort in reverse order (newest first)
	sort.Slice(appliedMigrations, func(i, j int) bool {
		return appliedMigrations[i].Version > appliedMigrations[j].Version
	})

	// Rollback the specified number of migrations
	count := 0
	for _, migration := range appliedMigrations {
		if count >= steps {
			break
		}
		if err := mm.RollbackMigration(migration); err != nil {
			return err
		}
		count++
	}

	logger.Infof("Rolled back %d migrations", count)
	return nil
}

// GetMigrationStatus returns the status of all migrations
func (mm *MigrationManager) GetMigrationStatus() ([]MigrationStatus, error) {
	applied, err := mm.GetAppliedMigrations()
	if err != nil {
		return nil, err
	}

	var status []MigrationStatus
	for _, migration := range mm.migrations {
		status = append(status, MigrationStatus{
			Version:   migration.Version,
			Name:      migration.Name,
			Applied:   applied[migration.Version],
			CreatedAt: migration.CreatedAt,
		})
	}

	return status, nil
}

// MigrationStatus represents the status of a migration
type MigrationStatus struct {
	Version   int64
	Name      string
	Applied   bool
	CreatedAt time.Time
}
