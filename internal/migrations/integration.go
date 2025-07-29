package migrations

import (
	"database/sql"
	"tickets/internal/logger"
)

// RunMigrationsOnStartup runs all pending migrations when the application starts
func RunMigrationsOnStartup(db *sql.DB, migrationsPath string) error {
	logger.Info("Running database migrations...")

	manager := NewMigrationManager(db)

	if err := manager.LoadMigrations(migrationsPath); err != nil {
		return err
	}

	if err := manager.MigrateUp(); err != nil {
		return err
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

// GetMigrationStatusOnStartup returns the current migration status
func GetMigrationStatusOnStartup(db *sql.DB, migrationsPath string) ([]MigrationStatus, error) {
	manager := NewMigrationManager(db)

	if err := manager.LoadMigrations(migrationsPath); err != nil {
		return nil, err
	}

	return manager.GetMigrationStatus()
}
