package repository

import (
	models "tickets/internal/models/domain"

	"github.com/jmoiron/sqlx"
)

// OrderRepository handles order and order item-related database operations
type OrderRepository struct {
	*BaseRepository
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(base *BaseRepository) *OrderRepository {
	return &OrderRepository{BaseRepository: base}
}

// CreateOrder creates a new order in the database
func (r *OrderRepository) CreateOrder(tx *sqlx.Tx, order *models.Order) error {
	query := `
		INSERT INTO orders (status, total_price) 
		VALUES ($1, $2) 
		RETURNING id, created_at, status, total_price`
	var createdAt int64
	err := tx.QueryRow(query, order.Status, order.TotalPrice).Scan(
		&order.ID, &createdAt, &order.Status, &order.TotalPrice)
	if err != nil {
		return err
	}
	order.CreatedAt = createdAt

	return nil
}
