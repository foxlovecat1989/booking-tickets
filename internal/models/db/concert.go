package db

import (
	models "tickets/internal/models/domain"

	"github.com/shopspring/decimal"
)

type Concert struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Location    string `db:"location"`
	Description string `db:"description"`
	CreatedAt   int64  `db:"created_at"`
}

func (c *Concert) ToConcert() *models.Concert {
	return &models.Concert{
		ID:          c.ID,
		Name:        c.Name,
		Location:    c.Location,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
	}
}

type ConcertSession struct {
	ID            int             `db:"id"`
	ConcertID     int             `db:"concert_id"`
	StartTime     int64           `db:"start_time"`
	EndTime       int64           `db:"end_time"`
	Venue         string          `db:"venue"`
	NumberOfSeats int             `db:"number_of_seats"`
	Price         decimal.Decimal `db:"price"`
}

func (c *ConcertSession) ToConcertSession() *models.ConcertSession {
	return &models.ConcertSession{
		ID:        c.ID,
		ConcertID: c.ConcertID,
		StartTime: c.StartTime,
		EndTime:   c.EndTime,
		Venue:     c.Venue,
		Price:     c.Price,
	}
}
