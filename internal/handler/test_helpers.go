package handler

import (
	"testing"

	"tickets/internal/repository"
	"tickets/internal/service"
)

// SetupTestHandler creates a test handler with a test database
func SetupTestHandler(t *testing.T) (*GRPCHandler, func()) {
	baseRepo, cleanup := repository.SetupTestDB(t)

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	return handler, cleanup
}

// SetupTestHandlerWithData creates a test handler with test data
func SetupTestHandlerWithData(t *testing.T) (*GRPCHandler, func()) {
	baseRepo, cleanup := repository.SetupTestDB(t)

	// Insert test data
	err := insertTestData(baseRepo)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	return handler, cleanup
}

// insertTestData inserts test data into the database
func insertTestData(baseRepo *repository.BaseRepository) error {
	// Insert test concert
	concertQuery := `
		INSERT INTO concerts (name, location, description) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	var concertID int
	err := baseRepo.GetDB().QueryRow(concertQuery,
		"Test Concert",
		"Test Venue",
		"Test Description").Scan(&concertID)
	if err != nil {
		return err
	}

	// Insert test concert session
	sessionQuery := `
		INSERT INTO concert_sessions (concert_id, start_time, end_time, venue, number_of_seats, price) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id`

	var sessionID int
	err = baseRepo.GetDB().QueryRow(sessionQuery,
		concertID,
		1735689600000, // Dec 31, 2024 8:00 PM
		1735700400000, // Dec 31, 2024 11:00 PM
		"Test Arena",
		100,
		99.99).Scan(&sessionID)
	if err != nil {
		return err
	}

	// Insert test tickets
	ticketQuery := `
		INSERT INTO tickets (session_id, status) 
		VALUES ($1, $2)`

	for i := 0; i < 10; i++ {
		_, err = baseRepo.GetDB().Exec(ticketQuery, sessionID, "available")
		if err != nil {
			return err
		}
	}

	return nil
}
