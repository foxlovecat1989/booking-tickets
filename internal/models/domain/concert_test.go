package models

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestConcert_Validation(t *testing.T) {
	tests := []struct {
		name    string
		concert Concert
		isValid bool
	}{
		{
			name: "valid concert",
			concert: Concert{
				Name:        "Rock Concert 2024",
				Location:    "Madison Square Garden",
				Description: "The biggest rock concert of the year",
				CreatedAt:   time.Now().UnixMilli(),
			},
			isValid: true,
		},
		{
			name: "missing name",
			concert: Concert{
				Location:    "Madison Square Garden",
				Description: "The biggest rock concert of the year",
				CreatedAt:   time.Now().UnixMilli(),
			},
			isValid: false,
		},
		{
			name: "missing location",
			concert: Concert{
				Name:        "Rock Concert 2024",
				Description: "The biggest rock concert of the year",
				CreatedAt:   time.Now().UnixMilli(),
			},
			isValid: false,
		},
		{
			name: "empty name",
			concert: Concert{
				Name:        "",
				Location:    "Madison Square Garden",
				Description: "The biggest rock concert of the year",
				CreatedAt:   time.Now().UnixMilli(),
			},
			isValid: false,
		},
		{
			name: "empty location",
			concert: Concert{
				Name:        "Rock Concert 2024",
				Location:    "",
				Description: "The biggest rock concert of the year",
				CreatedAt:   time.Now().UnixMilli(),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validation library
			// like go-playground/validator to validate struct tags
			isValid := tt.concert.Name != "" && tt.concert.Location != ""
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestConcertSession_Validation(t *testing.T) {
	now := time.Now()
	startTime := now.Add(time.Hour).UnixMilli()
	endTime := now.Add(2 * time.Hour).UnixMilli()

	tests := []struct {
		name    string
		session ConcertSession
		isValid bool
	}{
		{
			name: "valid session",
			session: ConcertSession{
				ConcertID: 1,
				StartTime: startTime,
				EndTime:   endTime,
				Venue:     "Main Arena",
				Price:     decimal.NewFromFloat(99.99),
			},
			isValid: true,
		},
		{
			name: "missing concert_id",
			session: ConcertSession{
				StartTime: startTime,
				EndTime:   endTime,
				Venue:     "Main Arena",
				Price:     decimal.NewFromFloat(99.99),
			},
			isValid: false,
		},
		{
			name: "missing start_time",
			session: ConcertSession{
				ConcertID: 1,
				EndTime:   endTime,
				Venue:     "Main Arena",
				Price:     decimal.NewFromFloat(99.99),
			},
			isValid: false,
		},
		{
			name: "missing end_time",
			session: ConcertSession{
				ConcertID: 1,
				StartTime: startTime,
				Venue:     "Main Arena",
				Price:     decimal.NewFromFloat(99.99),
			},
			isValid: false,
		},
		{
			name: "missing venue",
			session: ConcertSession{
				ConcertID: 1,
				StartTime: startTime,
				EndTime:   endTime,
				Price:     decimal.NewFromFloat(99.99),
			},
			isValid: false,
		},
		{
			name: "end_time before start_time",
			session: ConcertSession{
				ConcertID: 1,
				StartTime: endTime,
				EndTime:   startTime,
				Venue:     "Main Arena",
				Price:     decimal.NewFromFloat(99.99),
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validation library
			isValid := tt.session.ConcertID > 0 &&
				tt.session.StartTime > 0 &&
				tt.session.EndTime > 0 &&
				tt.session.Venue != "" &&
				tt.session.EndTime > tt.session.StartTime &&
				tt.session.Price.GreaterThanOrEqual(decimal.Zero)
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestConcertSession_WithConcert(t *testing.T) {
	concert := &Concert{
		ID:          1,
		Name:        "Rock Concert 2024",
		Location:    "Madison Square Garden",
		Description: "The biggest rock concert of the year",
		CreatedAt:   time.Now().UnixMilli(),
	}

	session := ConcertSession{
		ID:        1,
		ConcertID: 1,
		StartTime: time.Now().Add(time.Hour).UnixMilli(),
		EndTime:   time.Now().Add(2 * time.Hour).UnixMilli(),
		Venue:     "Main Arena",
		Price:     decimal.NewFromFloat(99.99),
		Concert:   concert,
	}

	assert.NotNil(t, session.Concert)
	assert.Equal(t, session.ID, 1)
	assert.Equal(t, session.ConcertID, 1)
	assert.Greater(t, session.StartTime, int64(0))
	assert.Greater(t, session.EndTime, session.StartTime)
	assert.Equal(t, "Main Arena", session.Venue)
	assert.Equal(t, decimal.NewFromFloat(99.99), session.Price)
	assert.Equal(t, concert.ID, session.Concert.ID)
	assert.Equal(t, concert.Name, session.Concert.Name)
	assert.Equal(t, concert.Location, session.Concert.Location)
}

func TestConcertSession_PriceHandling(t *testing.T) {
	tests := []struct {
		name     string
		price    decimal.Decimal
		expected string
	}{
		{
			name:     "zero price",
			price:    decimal.Zero,
			expected: "0",
		},
		{
			name:     "positive price",
			price:    decimal.NewFromFloat(99.99),
			expected: "99.99",
		},
		{
			name:     "large price",
			price:    decimal.NewFromFloat(999999.99),
			expected: "999999.99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := ConcertSession{
				Price: tt.price,
			}

			assert.Equal(t, tt.expected, session.Price.String())
		})
	}
}
