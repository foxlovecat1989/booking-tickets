package models

import "github.com/shopspring/decimal"

// Concert represents a concert in the system
type Concert struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
}

// ConcertSession represents a concert session
type ConcertSession struct {
	ID        int             `json:"id"`
	ConcertID int             `json:"concert_id" binding:"required"`
	StartTime int64           `json:"start_time" binding:"required"`
	EndTime   int64           `json:"end_time" binding:"required"`
	Venue     string          `json:"venue" binding:"required"`
	Price     decimal.Decimal `json:"price"`
	Concert   *Concert        `json:"concert,omitempty"`
}
