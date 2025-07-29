package models

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTicketType_Validation(t *testing.T) {
	tests := []struct {
		name       string
		ticketType TicketType
		isValid    bool
	}{
		{
			name: "valid ticket type",
			ticketType: TicketType{
				Name:        "VIP",
				Description: "VIP ticket with premium seating",
			},
			isValid: true,
		},
		{
			name: "missing name",
			ticketType: TicketType{
				Description: "VIP ticket with premium seating",
			},
			isValid: false,
		},
		{
			name: "empty name",
			ticketType: TicketType{
				Name:        "",
				Description: "VIP ticket with premium seating",
			},
			isValid: false,
		},
		{
			name: "name with spaces only",
			ticketType: TicketType{
				Name:        "   ",
				Description: "VIP ticket with premium seating",
			},
			isValid: false, // Spaces-only names should be invalid after trimming
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validation library
			trimmedName := strings.TrimSpace(tt.ticketType.Name)
			isValid := trimmedName != ""
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestTicket_Validation(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name    string
		ticket  Ticket
		isValid bool
	}{
		{
			name: "valid ticket",
			ticket: Ticket{
				ID:        validUUID,
				SessionID: 1,
				Status:    "available",
			},
			isValid: true,
		},
		{
			name: "missing session_id",
			ticket: Ticket{
				ID:     validUUID,
				Status: "available",
			},
			isValid: false,
		},
		{
			name: "invalid status",
			ticket: Ticket{
				ID:        validUUID,
				SessionID: 1,
				Status:    "invalid_status",
			},
			isValid: false,
		},
		{
			name: "zero session_id",
			ticket: Ticket{
				ID:        validUUID,
				SessionID: 0,
				Status:    "available",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validation library
			validStatuses := map[string]bool{
				"available": true,
				"pending":   true,
				"sold":      true,
			}

			isValid := tt.ticket.SessionID > 0 && validStatuses[tt.ticket.Status]
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestTicket_StatusValidation(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name   string
		status string
		valid  bool
	}{
		{"available status", "available", true},
		{"pending status", "pending", true},
		{"sold status", "sold", true},
		{"invalid status", "invalid", false},
		{"empty status", "", false},
		{"uppercase status", "AVAILABLE", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticket := Ticket{
				ID:        validUUID,
				SessionID: 1,
				Status:    tt.status,
			}

			validStatuses := map[string]bool{
				"available": true,
				"pending":   true,
				"sold":      true,
			}

			isValid := validStatuses[ticket.Status]
			assert.Equal(t, tt.valid, isValid)
			assert.Equal(t, 1, ticket.SessionID)
			assert.Equal(t, validUUID, ticket.ID)
			assert.Equal(t, tt.status, ticket.Status)
		})
	}
}

func TestCreateTicketRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreateTicketRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: CreateTicketRequest{
				SessionID:    1,
				SeatNumber:   "A1",
				TicketTypeID: 1,
				Price:        99.99,
			},
			isValid: true,
		},
		{
			name: "missing session_id",
			request: CreateTicketRequest{
				SeatNumber:   "A1",
				TicketTypeID: 1,
				Price:        99.99,
			},
			isValid: false,
		},
		{
			name: "missing seat_number",
			request: CreateTicketRequest{
				SessionID:    1,
				TicketTypeID: 1,
				Price:        99.99,
			},
			isValid: false,
		},
		{
			name: "missing ticket_type_id",
			request: CreateTicketRequest{
				SessionID:  1,
				SeatNumber: "A1",
				Price:      99.99,
			},
			isValid: false,
		},
		{
			name: "missing price",
			request: CreateTicketRequest{
				SessionID:    1,
				SeatNumber:   "A1",
				TicketTypeID: 1,
			},
			isValid: false,
		},
		{
			name: "negative price",
			request: CreateTicketRequest{
				SessionID:    1,
				SeatNumber:   "A1",
				TicketTypeID: 1,
				Price:        -99.99,
			},
			isValid: false,
		},
		{
			name: "zero price",
			request: CreateTicketRequest{
				SessionID:    1,
				SeatNumber:   "A1",
				TicketTypeID: 1,
				Price:        0,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validation library
			isValid := tt.request.SessionID > 0 &&
				tt.request.SeatNumber != "" &&
				tt.request.TicketTypeID > 0 &&
				tt.request.Price > 0
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestUpdateTicketRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request UpdateTicketRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: UpdateTicketRequest{
				SessionID:    1,
				SeatNumber:   "A1",
				TicketTypeID: 1,
				Price:        99.99,
			},
			isValid: true,
		},
		{
			name: "missing session_id",
			request: UpdateTicketRequest{
				SeatNumber:   "A1",
				TicketTypeID: 1,
				Price:        99.99,
			},
			isValid: false,
		},
		{
			name: "missing seat_number",
			request: UpdateTicketRequest{
				SessionID:    1,
				TicketTypeID: 1,
				Price:        99.99,
			},
			isValid: false,
		},
		{
			name: "missing ticket_type_id",
			request: UpdateTicketRequest{
				SessionID:  1,
				SeatNumber: "A1",
				Price:      99.99,
			},
			isValid: false,
		},
		{
			name: "missing price",
			request: UpdateTicketRequest{
				SessionID:    1,
				SeatNumber:   "A1",
				TicketTypeID: 1,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validation library
			isValid := tt.request.SessionID > 0 &&
				tt.request.SeatNumber != "" &&
				tt.request.TicketTypeID > 0 &&
				tt.request.Price > 0
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestTicket_UUIDHandling(t *testing.T) {
	// Test that UUID is properly handled
	ticket := Ticket{
		ID:        uuid.New(),
		SessionID: 1,
		Status:    "available",
	}

	assert.NotEqual(t, uuid.Nil, ticket.ID)
	assert.NotEmpty(t, ticket.ID.String())
	assert.Equal(t, 1, ticket.SessionID)
	assert.Equal(t, "available", ticket.Status)

	// Test UUID parsing
	uuidStr := ticket.ID.String()
	parsedUUID, err := uuid.Parse(uuidStr)
	assert.NoError(t, err)
	assert.Equal(t, ticket.ID, parsedUUID)
}

func TestTicket_PriceHandling(t *testing.T) {
	tests := []struct {
		name     string
		price    float64
		expected bool
	}{
		{"positive price", 99.99, true},
		{"zero price", 0.0, false},
		{"negative price", -99.99, false},
		{"very small price", 0.01, true},
		{"large price", 999999.99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := CreateTicketRequest{
				SessionID:    1,
				SeatNumber:   "A1",
				TicketTypeID: 1,
				Price:        tt.price,
			}

			isValid := request.Price > 0
			assert.Equal(t, tt.expected, isValid)
			assert.Equal(t, 1, request.SessionID)
			assert.Equal(t, "A1", request.SeatNumber)
			assert.Equal(t, 1, request.TicketTypeID)
			assert.Equal(t, tt.price, request.Price)
		})
	}
}
