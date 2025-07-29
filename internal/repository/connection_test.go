package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabaseConnection(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	// Test that we can ping the database
	err := baseRepo.db.Ping()
	assert.NoError(t, err)

	// Test that we can execute a simple query
	var result int
	err = baseRepo.db.Get(&result, "SELECT 1")
	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

func TestDatabaseTransaction(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	// Test that we can start and commit a transaction
	tx, err := baseRepo.db.Beginx()
	require.NoError(t, err)
	defer func() {
		if err := tx.Rollback(); err != nil {
			t.Logf("Error rolling back transaction: %v", err)
		}
	}()

	// Test a simple query within transaction
	var result int
	err = tx.Get(&result, "SELECT 1")
	assert.NoError(t, err)
	assert.Equal(t, 1, result)

	// Commit the transaction
	err = tx.Commit()
	assert.NoError(t, err)
}
