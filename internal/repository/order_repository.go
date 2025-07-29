package repository

// OrderRepository handles order and order item-related database operations
type OrderRepository struct {
	*BaseRepository
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(base *BaseRepository) *OrderRepository {
	return &OrderRepository{BaseRepository: base}
}
