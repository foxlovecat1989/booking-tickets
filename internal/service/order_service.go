package service

import (
	"errors"
	models "tickets/internal/models/domain"
	"tickets/internal/repository"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// OrderService handles order-related business logic
type OrderService struct {
	orderRepo          *repository.OrderRepository
	concertSessionRepo *repository.ConcertSessionRepository
	ticketRepo         *repository.TicketRepository
}

// NewOrderService creates a new order service
func NewOrderService(base *BaseService) *OrderService {
	baseRepo := base.GetBaseRepository()
	return &OrderService{
		orderRepo:          repository.NewOrderRepository(baseRepo),
		concertSessionRepo: repository.NewConcertSessionRepository(baseRepo),
		ticketRepo:         repository.NewTicketRepository(baseRepo),
	}
}

// CreateOrderRequest represents the request structure for creating an order
type CreateOrderRequest struct {
	UserID           int `json:"user_id" binding:"required"`
	ConcertSessionID int `json:"concert_session_id" binding:"required"`
	NumberOfTickets  int `json:"number_of_tickets" binding:"required"`
}

// CreateOrderResponse represents the response structure for creating an order
type CreateOrderResponse struct {
	OrderID    int             `json:"order_id"`
	Status     string          `json:"status"`
	TicketIDs  []string        `json:"ticket_ids"`
	TotalPrice decimal.Decimal `json:"total_price"`
	CreatedAt  int64           `json:"created_at"`
}

// CreateOrder creates a new order
func (s *OrderService) CreateOrder(req *CreateOrderRequest) (*CreateOrderResponse, error) {
	// Validate request is not nil
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// Validate number of tickets is within valid range
	if req.NumberOfTickets <= 0 {
		return nil, errors.New("number of tickets must be greater than 0")
	}
	if req.NumberOfTickets > 3 {
		return nil, errors.New("maximum 3 tickets allowed per order")
	}

	var order *models.Order
	var tickets []models.Ticket

	// Execute everything in a transaction
	err := s.orderRepo.BaseRepository.WithTransaction(func(tx *sqlx.Tx) error {
		// Validate concert session exists
		concertSession, err := s.concertSessionRepo.GetConcertSessionByID(req.ConcertSessionID)
		if err != nil {
			return err
		}
		if concertSession == nil {
			return errors.New("concert session not found")
		}

		// Validate tickets are available
		tickets, err = s.ticketRepo.GetAvailableTicketsBySessionID(req.ConcertSessionID, req.NumberOfTickets)
		if err != nil {
			return err
		}
		if len(tickets) == 0 {
			return errors.New("no tickets available")
		}

		// Create order with basic information
		order = &models.Order{
			Status:     "pending",
			TotalPrice: decimal.NewFromInt(int64(len(tickets))).Mul(concertSession.Price),
		}

		// Create order in database
		err = s.orderRepo.CreateOrder(tx, order)
		if err != nil {
			return err
		}

		// Update ticket statuses to 'pending'
		err = s.ticketRepo.UpdateTicketStatuses(tx, tickets, "pending")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	ticketIDs := make([]string, len(tickets))
	for i, ticket := range tickets {
		ticketIDs[i] = ticket.ID.String()
	}

	return &CreateOrderResponse{
		OrderID:    order.ID,
		Status:     order.Status,
		TicketIDs:  ticketIDs,
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt,
	}, nil
}
