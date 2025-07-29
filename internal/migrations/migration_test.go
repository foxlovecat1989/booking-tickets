package migrations

import (
	"database/sql"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestMigrationManager_LoadMigrations(t *testing.T) {
	// Create a temporary directory for test migrations
	tempDir, err := os.MkdirTemp("", "migrations_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test migration files
	testMigrations := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "001_test.up.sql",
			content:  `CREATE TABLE test (id SERIAL PRIMARY KEY);`,
			expected: true,
		},
		{
			name:     "001_test.down.sql",
			content:  `DROP TABLE test;`,
			expected: true,
		},
		{
			name:     "002_another.up.sql",
			content:  `CREATE TABLE another (id SERIAL PRIMARY KEY);`,
			expected: true,
		},
		{
			name:     "invalid.sql",
			content:  `SELECT * FROM test;`,
			expected: false,
		},
	}

	for _, tm := range testMigrations {
		filePath := filepath.Join(tempDir, tm.name)
		if err := os.WriteFile(filePath, []byte(tm.content), 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", tm.name, err)
		}
	}

	// Create a mock database connection (we'll use a nil connection for this test)
	var db *sql.DB
	manager := NewMigrationManager(db)

	// Load migrations
	if err := manager.LoadMigrations(tempDir); err != nil {
		t.Fatalf("Failed to load migrations: %v", err)
	}

	// Verify migrations were loaded correctly
	expectedCount := 2 // Only valid migrations should be loaded
	if len(manager.migrations) != expectedCount {
		t.Errorf("Expected %d migrations, got %d", expectedCount, len(manager.migrations))
	}

	// Check that migrations are sorted by version
	for i := 1; i < len(manager.migrations); i++ {
		if manager.migrations[i-1].Version >= manager.migrations[i].Version {
			t.Errorf("Migrations not sorted correctly: %d >= %d",
				manager.migrations[i-1].Version, manager.migrations[i].Version)
		}
	}
}

func TestMigrationManager_ParseMigrationFileName(t *testing.T) {
	testCases := []struct {
		filename    string
		expectValid bool
		version     int64
		name        string
		direction   string
	}{
		{"001_test.up.sql", true, 1, "test", "up"},
		{"002_another_migration.down.sql", true, 2, "another_migration", "down"},
		{"invalid.sql", false, 0, "", ""},
		{"001_test.invalid.sql", true, 1, "test", "invalid"},
		{"abc_test.up.sql", false, 0, "", ""},
	}

	for _, tc := range testCases {
		parts := parseMigrationFileName(tc.filename)
		if tc.expectValid {
			if parts == nil {
				t.Errorf("Expected valid migration file %s, but got nil", tc.filename)
				continue
			}
			if parts.version != tc.version {
				t.Errorf("Expected version %d for %s, got %d", tc.version, tc.filename, parts.version)
			}
			if parts.name != tc.name {
				t.Errorf("Expected name %s for %s, got %s", tc.name, tc.filename, parts.name)
			}
			if parts.direction != tc.direction {
				t.Errorf("Expected direction %s for %s, got %s", tc.direction, tc.filename, parts.direction)
			}
		} else {
			if parts != nil {
				t.Errorf("Expected invalid migration file %s, but got valid parts", tc.filename)
			}
		}
	}
}

// Helper function to parse migration file name
type migrationParts struct {
	version   int64
	name      string
	direction string
}

func parseMigrationFileName(filename string) *migrationParts {
	// Parse migration file name: 001_initial_schema.up.sql
	parts := strings.Split(strings.TrimSuffix(filename, ".sql"), ".")
	if len(parts) != 2 {
		return nil
	}

	versionName := parts[0]
	direction := parts[1]

	// Extract version number
	versionParts := strings.SplitN(versionName, "_", 2)
	if len(versionParts) != 2 {
		return nil
	}

	version, err := strconv.ParseInt(versionParts[0], 10, 64)
	if err != nil {
		return nil
	}

	name := versionParts[1]

	return &migrationParts{
		version:   version,
		name:      name,
		direction: direction,
	}
}

func TestMigrationStatus_Formatting(t *testing.T) {
	status := MigrationStatus{
		Version:   1,
		Name:      "test_migration",
		Applied:   true,
		CreatedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
	}

	// Test that the struct can be created and accessed
	if status.Version != 1 {
		t.Errorf("Expected version 1, got %d", status.Version)
	}
	if status.Name != "test_migration" {
		t.Errorf("Expected name 'test_migration', got %s", status.Name)
	}
	if !status.Applied {
		t.Error("Expected applied to be true")
	}
	if status.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}
