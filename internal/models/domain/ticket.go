package models

import (
	"github.com/google/uuid"
)

// TicketType represents a ticket type
type TicketType struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// Ticket represents a ticket in the system
type Ticket struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SessionID int       `json:"session_id" db:"session_id"`
	Status    string    `json:"status" db:"status"`
}

// CreateTicketRequest represents the request structure for creating a ticket
type CreateTicketRequest struct {
	SessionID    int     `json:"session_id" binding:"required"`
	SeatNumber   string  `json:"seat_number" binding:"required"`
	TicketTypeID int     `json:"ticket_type_id" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
}

// UpdateTicketRequest represents the request structure for updating a ticket
type UpdateTicketRequest struct {
	SessionID    int     `json:"session_id" binding:"required"`
	SeatNumber   string  `json:"seat_number" binding:"required"`
	TicketTypeID int     `json:"ticket_type_id" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
}
