package repository

import (
	"database/sql"
	"tickets/internal/models/db"
	models "tickets/internal/models/domain"
)

// ConcertSessionRepository handles concert session-related database operations
type ConcertSessionRepository struct {
	*BaseRepository
}

// NewConcertSessionRepository creates a new concert session repository
func NewConcertSessionRepository(base *BaseRepository) *ConcertSessionRepository {
	return &ConcertSessionRepository{BaseRepository: base}
}

// GetConcertSessionByID retrieves a concert session by ID
func (r *ConcertSessionRepository) GetConcertSessionByID(id int) (*models.ConcertSession, error) {
	query := `SELECT id, concert_id, start_time, end_time, venue, number_of_seats, price FROM concert_sessions WHERE id = $1`

	var dbSession db.ConcertSession
	err := r.db.Get(&dbSession, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return dbSession.ToConcertSession(), nil
}
