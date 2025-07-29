package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Order represents an order in the system
type Order struct {
	ID         int             `json:"id"`
	CreatedAt  int64           `json:"created_at"`
	Status     string          `json:"status"`
	TotalPrice decimal.Decimal `json:"total_price"`
	Items      []OrderItem     `json:"items,omitempty"`
}

// OrderItem represents an order item
type OrderItem struct {
	ID       int             `json:"id"`
	OrderID  int             `json:"order_id" binding:"required"`
	TicketID uuid.UUID       `json:"ticket_id" binding:"required"`
	Price    decimal.Decimal `json:"price" binding:"required"`
	Ticket   *Ticket         `json:"ticket,omitempty"`
}
