package repository

import (
	"github.com/jmoiron/sqlx"
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

// UpdateTicketStatuses updates the status of multiple tickets
func (r *TicketRepository) UpdateTicketStatuses(tx *sqlx.Tx, tickets []models.Ticket, status string) error {
	query := `
	UPDATE tickets 
	SET status = $1 
	WHERE id = $2`

	for _, ticket := range tickets {
		_, err := tx.Exec(query, status, ticket.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
