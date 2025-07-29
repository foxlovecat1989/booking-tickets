package repository

import (
	models "tickets/internal/models/domain"
)

// TicketRepository handles ticket-related database operations
type TicketRepository struct {
	*BaseRepository
}

// NewTicketRepository creates a new ticket repository
func NewTicketRepository(base *BaseRepository) *TicketRepository {
	return &TicketRepository{BaseRepository: base}
}

// GetAvailableTicketsBySessionID retrieves available tickets for a session
func (r *TicketRepository) GetAvailableTicketsBySessionID(sessionID int, numberOfTickets int) ([]models.Ticket, error) {
	query := `
	SELECT id, session_id, status 
	FROM tickets 
	WHERE session_id = $1 AND status = 'available'
	ORDER BY id ASC
	LIMIT $2
	FOR UPDATE
	`

	var tickets []models.Ticket
	err := r.GetDB().Select(&tickets, query, sessionID, numberOfTickets)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}
