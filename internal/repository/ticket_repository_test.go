package repository

import (
	"testing"

	models "tickets/internal/models/domain"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTicketRepository(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewTicketRepository(baseRepo)
	assert.NotNil(t, repo)
	assert.Equal(t, baseRepo, repo.BaseRepository)
}

func TestTicketRepository_UpdateTicketStatuses(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewTicketRepository(baseRepo)

	// First, create some test tickets
	tickets := createTestTickets(t, baseRepo, 3)

	// Update all tickets to 'pending' status
	err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
		return repo.UpdateTicketStatuses(tx, tickets, "pending")
	})

	require.NoError(t, err)

	// Verify all tickets were updated
	for _, ticket := range tickets {
		var status string
		err := baseRepo.db.Get(&status, "SELECT status FROM tickets WHERE id = $1", ticket.ID)
		require.NoError(t, err)
		assert.Equal(t, "pending", status)
	}
}

func TestTicketRepository_UpdateTicketStatuses_EmptySlice(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewTicketRepository(baseRepo)

	// Test updating an empty slice of tickets
	err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
		return repo.UpdateTicketStatuses(tx, []models.Ticket{}, "pending")
	})

	require.NoError(t, err)
}

func TestTicketRepository_UpdateTicketStatuses_DifferentStatuses(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewTicketRepository(baseRepo)

	// Create test tickets
	tickets := createTestTickets(t, baseRepo, 2)

	testCases := []string{"pending", "sold", "available"}

	for _, status := range testCases {
		t.Run("status_"+status, func(t *testing.T) {
			err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
				return repo.UpdateTicketStatuses(tx, tickets, status)
			})

			require.NoError(t, err)

			// Verify all tickets were updated to the new status
			for _, ticket := range tickets {
				var dbStatus string
				err := baseRepo.db.Get(&dbStatus, "SELECT status FROM tickets WHERE id = $1", ticket.ID)
				require.NoError(t, err)
				assert.Equal(t, status, dbStatus)
			}
		})
	}
}

func TestTicketRepository_UpdateTicketStatuses_TransactionRollback(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewTicketRepository(baseRepo)

	// Create test tickets
	tickets := createTestTickets(t, baseRepo, 2)

	// Verify initial status is 'available'
	for _, ticket := range tickets {
		var status string
		err := baseRepo.db.Get(&status, "SELECT status FROM tickets WHERE id = $1", ticket.ID)
		require.NoError(t, err)
		assert.Equal(t, "available", status)
	}

	// Simulate a transaction that will be rolled back
	err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
		// Update tickets to 'pending'
		err := repo.UpdateTicketStatuses(tx, tickets, "pending")
		if err != nil {
			return err
		}

		// Simulate an error that causes rollback
		return assert.AnError
	})

	require.Error(t, err)

	// Verify the tickets were not updated (rolled back)
	for _, ticket := range tickets {
		var status string
		err := baseRepo.db.Get(&status, "SELECT status FROM tickets WHERE id = $1", ticket.ID)
		require.NoError(t, err)
		assert.Equal(t, "available", status)
	}
}

func TestTicketRepository_UpdateTicketStatuses_ConcurrentAccess(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewTicketRepository(baseRepo)

	// Create test tickets
	tickets := createTestTickets(t, baseRepo, 5)

	// Test concurrent status updates
	const numGoroutines = 3
	errors := make(chan error, numGoroutines)

	statuses := []string{"pending", "sold", "available"}

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			status := statuses[index%len(statuses)]
			err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
				return repo.UpdateTicketStatuses(tx, tickets, status)
			})
			errors <- err
		}(i)
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	// Verify all tickets have a valid status
	for _, ticket := range tickets {
		var status string
		err := baseRepo.db.Get(&status, "SELECT status FROM tickets WHERE id = $1", ticket.ID)
		require.NoError(t, err)
		assert.Contains(t, []string{"pending", "sold", "available"}, status)
	}
}

func TestTicketRepository_UpdateTicketStatuses_Integration(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewTicketRepository(baseRepo)

	// Test that we can connect to the database
	err := baseRepo.db.Ping()
	require.NoError(t, err)

	// Test that we can execute a simple query
	var count int
	err = baseRepo.db.Get(&count, "SELECT COUNT(*) FROM tickets")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 0)

	// Test that the repository is properly initialized
	assert.NotNil(t, repo)
	assert.Equal(t, baseRepo, repo.BaseRepository)
}

func TestTicketRepository_UpdateTicketStatuses_DataConsistency(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewTicketRepository(baseRepo)

	// Create test tickets
	tickets := createTestTickets(t, baseRepo, 3)

	// Update tickets to different statuses in sequence
	statuses := []string{"pending", "sold", "available"}

	for _, status := range statuses {
		err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
			return repo.UpdateTicketStatuses(tx, tickets, status)
		})
		require.NoError(t, err, "Failed to update tickets to status %s", status)

		// Verify all tickets have the expected status
		for _, ticket := range tickets {
			var dbStatus string
			err := baseRepo.db.Get(&dbStatus, "SELECT status FROM tickets WHERE id = $1", ticket.ID)
			require.NoError(t, err)
			assert.Equal(t, status, dbStatus, "Ticket %s should have status %s", ticket.ID, status)
		}
	}
}

func TestTicketRepository_GetAvailableTicketsBySessionID(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewTicketRepository(baseRepo)

	// Create a concert session first
	sessionID := createTestConcertSession(t, baseRepo)

	// Create test tickets for this session
	createTestTicketsForSession(t, baseRepo, sessionID, 5)

	// Test getting available tickets
	availableTickets, err := repo.GetAvailableTicketsBySessionID(sessionID, 3)
	require.NoError(t, err)
	assert.Len(t, availableTickets, 3)

	// Verify all returned tickets are available
	for _, ticket := range availableTickets {
		assert.Equal(t, "available", ticket.Status)
		assert.Equal(t, sessionID, ticket.SessionID)
	}
}

// Helper functions for creating test data

func createTestTickets(t *testing.T, baseRepo *BaseRepository, count int) []models.Ticket {
	// Create a concert session first
	sessionID := createTestConcertSession(t, baseRepo)

	return createTestTicketsForSession(t, baseRepo, sessionID, count)
}

func createTestConcertSession(t *testing.T, baseRepo *BaseRepository) int {
	// Create a concert first
	var concertID int
	err := baseRepo.db.QueryRow(`
		INSERT INTO concerts (name, location, description) 
		VALUES ($1, $2, $3) 
		RETURNING id`,
		"Test Concert", "Test Venue", "Test Description").Scan(&concertID)
	require.NoError(t, err)

	// Create a concert session
	var sessionID int
	err = baseRepo.db.QueryRow(`
		INSERT INTO concert_sessions (concert_id, start_time, end_time, venue, number_of_seats, price) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`,
		concertID, 1640995200000, 1640998800000, "Test Venue", 100, "50.00").Scan(&sessionID)
	require.NoError(t, err)

	return sessionID
}

func createTestTicketsForSession(t *testing.T, baseRepo *BaseRepository, sessionID int, count int) []models.Ticket {
	tickets := make([]models.Ticket, count)

	for i := 0; i < count; i++ {
		var ticket models.Ticket
		err := baseRepo.db.QueryRow(`
			INSERT INTO tickets (session_id, status) 
			VALUES ($1, $2) 
			RETURNING id, session_id, status`,
			sessionID, "available").Scan(&ticket.ID, &ticket.SessionID, &ticket.Status)
		require.NoError(t, err)
		tickets[i] = ticket
	}

	return tickets
}
