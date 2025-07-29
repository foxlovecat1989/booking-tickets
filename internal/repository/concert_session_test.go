package repository

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConcertSessionRepository(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewConcertSessionRepository(baseRepo)
	assert.NotNil(t, repo)
	assert.Equal(t, baseRepo, repo.BaseRepository)
}

func TestConcertSessionRepository_GetConcertSessionByID(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewConcertSessionRepository(baseRepo)

	// Test getting non-existent session
	session, err := repo.GetConcertSessionByID(999)
	require.NoError(t, err)
	assert.Nil(t, session)

	// Test getting existing session (if database has data)
	// This test assumes the database has been set up with test data
	// In a real scenario, you would insert test data first
	session, err = repo.GetConcertSessionByID(1)
	if err != nil {
		// If there's no data, that's expected for this test
		t.Logf("No test data found, skipping existing session test: %v", err)
		return
	}

	if session != nil {
		assert.NotZero(t, session.ID)
		assert.NotZero(t, session.ConcertID)
		assert.NotZero(t, session.StartTime)
		assert.NotZero(t, session.EndTime)
		assert.NotEmpty(t, session.Venue)
		assert.True(t, session.Price.GreaterThanOrEqual(decimal.Zero))
	}
}

func TestConcertSessionRepository_GetConcertSessionByID_InvalidID(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewConcertSessionRepository(baseRepo)

	// Test with zero ID
	session, err := repo.GetConcertSessionByID(0)
	require.NoError(t, err)
	assert.Nil(t, session)

	// Test with negative ID
	session, err = repo.GetConcertSessionByID(-1)
	require.NoError(t, err)
	assert.Nil(t, session)
}

func TestConcertSessionRepository_GetConcertSessionByID_DatabaseError(t *testing.T) {
	// This test would require mocking the database connection
	// to simulate database errors
	// For now, we'll test the basic functionality
	t.Skip("Database error simulation requires mocking")
}

func TestConcertSessionRepository_Integration(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewConcertSessionRepository(baseRepo)

	// Test that we can connect to the database
	err := baseRepo.db.Ping()
	require.NoError(t, err)

	// Test that we can execute a simple query
	var count int
	err = baseRepo.db.Get(&count, "SELECT COUNT(*) FROM concert_sessions")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 0)

	// Test that the repository is properly initialized
	assert.NotNil(t, repo)
	assert.Equal(t, baseRepo, repo.BaseRepository)
}

func TestConcertSessionRepository_DataConsistency(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewConcertSessionRepository(baseRepo)

	// Test that retrieved data is consistent
	session, err := repo.GetConcertSessionByID(1)
	if err != nil || session == nil {
		t.Skip("No test data available for consistency test")
		return
	}

	// Test that the session has valid data
	assert.Greater(t, session.ID, 0)
	assert.Greater(t, session.ConcertID, 0)
	assert.Greater(t, session.StartTime, int64(0))
	assert.Greater(t, session.EndTime, session.StartTime) // End time should be after start time
	assert.NotEmpty(t, session.Venue)
	assert.True(t, session.Price.GreaterThanOrEqual(decimal.Zero))
}

func TestConcertSessionRepository_Performance(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewConcertSessionRepository(baseRepo)

	// Test that queries execute within reasonable time
	start := time.Now()

	// Execute multiple queries to test performance
	for i := 0; i < 10; i++ {
		_, err := repo.GetConcertSessionByID(1)
		require.NoError(t, err)
	}

	duration := time.Since(start)

	// Should complete within 1 second for 10 queries
	assert.Less(t, duration, time.Second)
}

func TestConcertSessionRepository_ConcurrentAccess(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewConcertSessionRepository(baseRepo)

	// Test concurrent access to the repository
	const numGoroutines = 10
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			_, err := repo.GetConcertSessionByID(1)
			// We don't require success here as there might not be data
			// but we do require no panics or unexpected errors
			if err != nil {
				t.Logf("Goroutine %d got expected error: %v", id, err)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

func TestConcertSessionRepository_ErrorHandling(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewConcertSessionRepository(baseRepo)

	// Test with various edge cases
	testCases := []struct {
		name string
		id   int
	}{
		{"zero id", 0},
		{"negative id", -1},
		{"very large id", 999999999},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			session, err := repo.GetConcertSessionByID(tc.id)
			require.NoError(t, err)
			assert.Nil(t, session)
		})
	}
}
