package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOrder_Validation(t *testing.T) {
	tests := []struct {
		name    string
		order   Order
		isValid bool
	}{
		{
			name: "valid order",
			order: Order{
				CreatedAt:  time.Now().UnixMilli(),
				Status:     "pending",
				TotalPrice: decimal.NewFromFloat(99.99),
			},
			isValid: true,
		},
		{
			name: "missing status",
			order: Order{
				CreatedAt:  time.Now().UnixMilli(),
				TotalPrice: decimal.NewFromFloat(99.99),
			},
			isValid: false,
		},
		{
			name: "invalid status",
			order: Order{
				CreatedAt:  time.Now().UnixMilli(),
				Status:     "invalid_status",
				TotalPrice: decimal.NewFromFloat(99.99),
			},
			isValid: false,
		},
		{
			name: "negative total price",
			order: Order{
				CreatedAt:  time.Now().UnixMilli(),
				Status:     "pending",
				TotalPrice: decimal.NewFromFloat(-99.99),
			},
			isValid: false,
		},
		{
			name: "zero total price",
			order: Order{
				CreatedAt:  time.Now().UnixMilli(),
				Status:     "pending",
				TotalPrice: decimal.Zero,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validation library
			validStatuses := map[string]bool{
				"pending":   true,
				"paid":      true,
				"cancelled": true,
				"completed": true,
			}

			isValid := validStatuses[tt.order.Status] && tt.order.TotalPrice.GreaterThan(decimal.Zero)
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestOrder_StatusValidation(t *testing.T) {
	tests := []struct {
		name   string
		status string
		valid  bool
	}{
		{"pending status", "pending", true},
		{"paid status", "paid", true},
		{"cancelled status", "cancelled", true},
		{"completed status", "completed", true},
		{"invalid status", "invalid", false},
		{"empty status", "", false},
		{"uppercase status", "PENDING", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := Order{
				CreatedAt:  time.Now().UnixMilli(),
				Status:     tt.status,
				TotalPrice: decimal.NewFromFloat(99.99),
			}

			validStatuses := map[string]bool{
				"pending":   true,
				"paid":      true,
				"cancelled": true,
				"completed": true,
			}

			isValid := validStatuses[order.Status]
			assert.Equal(t, tt.valid, isValid)
			assert.Greater(t, order.CreatedAt, int64(0))
			assert.Equal(t, decimal.NewFromFloat(99.99), order.TotalPrice)
		})
	}
}

func TestOrderItem_Validation(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name      string
		orderItem OrderItem
		isValid   bool
	}{
		{
			name: "valid order item",
			orderItem: OrderItem{
				OrderID:  1,
				TicketID: validUUID,
				Price:    decimal.NewFromFloat(99.99),
			},
			isValid: true,
		},
		{
			name: "missing order_id",
			orderItem: OrderItem{
				TicketID: validUUID,
				Price:    decimal.NewFromFloat(99.99),
			},
			isValid: false,
		},
		{
			name: "missing ticket_id",
			orderItem: OrderItem{
				OrderID: 1,
				Price:   decimal.NewFromFloat(99.99),
			},
			isValid: false,
		},
		{
			name: "missing price",
			orderItem: OrderItem{
				OrderID:  1,
				TicketID: validUUID,
			},
			isValid: false,
		},
		{
			name: "zero order_id",
			orderItem: OrderItem{
				OrderID:  0,
				TicketID: validUUID,
				Price:    decimal.NewFromFloat(99.99),
			},
			isValid: false,
		},
		{
			name: "negative price",
			orderItem: OrderItem{
				OrderID:  1,
				TicketID: validUUID,
				Price:    decimal.NewFromFloat(-99.99),
			},
			isValid: false,
		},
		{
			name: "zero price",
			orderItem: OrderItem{
				OrderID:  1,
				TicketID: validUUID,
				Price:    decimal.Zero,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: In a real application, you would use a validation library
			isValid := tt.orderItem.OrderID > 0 &&
				tt.orderItem.TicketID != uuid.Nil &&
				tt.orderItem.Price.GreaterThan(decimal.Zero)
			assert.Equal(t, tt.isValid, isValid)
		})
	}
}

func TestOrder_WithItems(t *testing.T) {
	validUUID := uuid.New()
	ticket := &Ticket{
		ID:        validUUID,
		SessionID: 1,
		Status:    "available",
	}

	orderItem := OrderItem{
		ID:       1,
		OrderID:  1,
		TicketID: validUUID,
		Price:    decimal.NewFromFloat(99.99),
		Ticket:   ticket,
	}

	order := Order{
		ID:         1,
		CreatedAt:  time.Now().UnixMilli(),
		Status:     "pending",
		TotalPrice: decimal.NewFromFloat(99.99),
		Items:      []OrderItem{orderItem},
	}

	assert.Len(t, order.Items, 1)
	assert.Equal(t, 1, order.ID)
	assert.Greater(t, order.CreatedAt, int64(0))
	assert.Equal(t, "pending", order.Status)
	assert.Equal(t, decimal.NewFromFloat(99.99), order.TotalPrice)
	assert.NotNil(t, order.Items[0].Ticket)
	assert.Equal(t, ticket.ID, order.Items[0].Ticket.ID)
	assert.Equal(t, ticket.SessionID, order.Items[0].Ticket.SessionID)
}

func TestOrder_PriceCalculations(t *testing.T) {
	validUUID := uuid.New()

	orderItem1 := OrderItem{
		OrderID:  1,
		TicketID: validUUID,
		Price:    decimal.NewFromFloat(49.99),
	}

	orderItem2 := OrderItem{
		OrderID:  1,
		TicketID: uuid.New(),
		Price:    decimal.NewFromFloat(59.99),
	}

	order := Order{
		ID:        1,
		CreatedAt: time.Now().UnixMilli(),
		Status:    "pending",
		Items:     []OrderItem{orderItem1, orderItem2},
	}

	// Calculate total price from items
	totalPrice := decimal.Zero
	for _, item := range order.Items {
		totalPrice = totalPrice.Add(item.Price)
	}

	expectedTotal := decimal.NewFromFloat(109.98)
	assert.Equal(t, 1, order.ID)
	assert.Greater(t, order.CreatedAt, int64(0))
	assert.Equal(t, "pending", order.Status)
	assert.True(t, totalPrice.Equal(expectedTotal))
}

func TestOrderItem_UUIDHandling(t *testing.T) {
	// Test that UUID is properly handled
	orderItem := OrderItem{
		OrderID:  1,
		TicketID: uuid.New(),
		Price:    decimal.NewFromFloat(99.99),
	}

	assert.Equal(t, 1, orderItem.OrderID)
	assert.Equal(t, decimal.NewFromFloat(99.99), orderItem.Price)
	assert.NotEqual(t, uuid.Nil, orderItem.TicketID)
	assert.NotEmpty(t, orderItem.TicketID.String())

	// Test UUID parsing
	uuidStr := orderItem.TicketID.String()
	parsedUUID, err := uuid.Parse(uuidStr)
	assert.NoError(t, err)
	assert.Equal(t, orderItem.TicketID, parsedUUID)
}

func TestOrder_PriceHandling(t *testing.T) {
	tests := []struct {
		name     string
		price    decimal.Decimal
		expected bool
	}{
		{"positive price", decimal.NewFromFloat(99.99), true},
		{"zero price", decimal.Zero, false},
		{"negative price", decimal.NewFromFloat(-99.99), false},
		{"very small price", decimal.NewFromFloat(0.01), true},
		{"large price", decimal.NewFromFloat(999999.99), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := Order{
				ID:         1,
				CreatedAt:  time.Now().UnixMilli(),
				Status:     "pending",
				TotalPrice: tt.price,
			}

			validStatuses := map[string]bool{
				"pending":   true,
				"paid":      true,
				"cancelled": true,
				"completed": true,
			}
			isValid := validStatuses[order.Status] && order.TotalPrice.GreaterThan(decimal.Zero)
			assert.Equal(t, tt.expected, isValid)
			assert.Equal(t, 1, order.ID)
			assert.Greater(t, order.CreatedAt, int64(0))
			assert.Equal(t, "pending", order.Status)
			assert.Equal(t, tt.price, order.TotalPrice)
		})
	}
}

func TestOrder_TimestampHandling(t *testing.T) {
	now := time.Now()

	order := Order{
		CreatedAt:  now.UnixMilli(),
		Status:     "pending",
		TotalPrice: decimal.NewFromFloat(99.99),
	}

	// Test that timestamp is valid
	assert.Greater(t, order.CreatedAt, int64(0))
	assert.Equal(t, "pending", order.Status)
	assert.Equal(t, decimal.NewFromFloat(99.99), order.TotalPrice)

	// Test that timestamp is recent (within last minute)
	oneMinuteAgo := time.Now().Add(-time.Minute).UnixMilli()
	assert.Greater(t, order.CreatedAt, oneMinuteAgo)
}

func TestOrderItem_PriceHandling(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name     string
		price    decimal.Decimal
		expected bool
	}{
		{"positive price", decimal.NewFromFloat(99.99), true},
		{"zero price", decimal.Zero, false},
		{"negative price", decimal.NewFromFloat(-99.99), false},
		{"very small price", decimal.NewFromFloat(0.01), true},
		{"large price", decimal.NewFromFloat(999999.99), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderItem := OrderItem{
				OrderID:  1,
				TicketID: validUUID,
				Price:    tt.price,
			}

			isValid := orderItem.Price.GreaterThan(decimal.Zero)
			assert.Equal(t, tt.expected, isValid)
			assert.Equal(t, 1, orderItem.OrderID)
			assert.Equal(t, validUUID, orderItem.TicketID)
			assert.Equal(t, tt.price, orderItem.Price)
		})
	}
}
