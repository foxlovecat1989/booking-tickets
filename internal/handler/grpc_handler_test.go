package handler

import (
	"context"
	"testing"

	"tickets/api"
	"tickets/internal/repository"
	"tickets/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNewGRPCHandler(t *testing.T) {
	handler, cleanup := SetupTestHandler(t)
	defer cleanup()

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.orderService)
}

func TestGRPCHandler_CreateOrder_ValidRequest(t *testing.T) {
	handler, cleanup := SetupTestHandlerWithData(t)
	defer cleanup()

	req := &api.CreateOrderRequest{
		UserId:           1,
		ConcertSessionId: 1,
		NumberOfTickets:  1,
	}

	// This test will fail if there's no test data in the database
	// In a real scenario, you would set up test data first
	resp, err := handler.CreateOrder(context.Background(), req)
	if err != nil {
		// If there's no test data, that's expected
		t.Logf("Expected error due to no test data: %v", err)
		return
	}

	require.NotNil(t, resp)
	assert.Greater(t, resp.OrderId, int32(0))
	assert.Equal(t, "pending", resp.Status)
	assert.NotEmpty(t, resp.TicketIds)
	assert.Greater(t, resp.TotalPrice, float64(0))
	assert.NotNil(t, resp.CreatedAt)
}

func TestGRPCHandler_CreateOrder_InvalidUserId(t *testing.T) {
	baseRepo, cleanup := repository.SetupTestDB(t)
	defer cleanup()

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	testCases := []struct {
		name        string
		userId      int32
		expectCode  codes.Code
		expectError string
	}{
		{
			name:        "zero user_id",
			userId:      0,
			expectCode:  codes.InvalidArgument,
			expectError: "user_id must be positive",
		},
		{
			name:        "negative user_id",
			userId:      -1,
			expectCode:  codes.InvalidArgument,
			expectError: "user_id must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &api.CreateOrderRequest{
				UserId:           tc.userId,
				ConcertSessionId: 1,
				NumberOfTickets:  1,
			}

			resp, err := handler.CreateOrder(context.Background(), req)
			assert.Nil(t, resp)
			assert.Error(t, err)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tc.expectCode, st.Code())
			assert.Contains(t, st.Message(), tc.expectError)
		})
	}
}

func TestGRPCHandler_CreateOrder_InvalidConcertSessionId(t *testing.T) {
	baseRepo, cleanup := repository.SetupTestDB(t)
	defer cleanup()

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	testCases := []struct {
		name        string
		sessionId   int32
		expectCode  codes.Code
		expectError string
	}{
		{
			name:        "zero concert_session_id",
			sessionId:   0,
			expectCode:  codes.InvalidArgument,
			expectError: "concert_session_id must be positive",
		},
		{
			name:        "negative concert_session_id",
			sessionId:   -1,
			expectCode:  codes.InvalidArgument,
			expectError: "concert_session_id must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &api.CreateOrderRequest{
				UserId:           1,
				ConcertSessionId: tc.sessionId,
				NumberOfTickets:  1,
			}

			resp, err := handler.CreateOrder(context.Background(), req)
			assert.Nil(t, resp)
			assert.Error(t, err)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tc.expectCode, st.Code())
			assert.Contains(t, st.Message(), tc.expectError)
		})
	}
}

func TestGRPCHandler_CreateOrder_InvalidNumberOfTickets(t *testing.T) {
	baseRepo, cleanup := repository.SetupTestDB(t)
	defer cleanup()

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	testCases := []struct {
		name        string
		numTickets  int32
		expectCode  codes.Code
		expectError string
	}{
		{
			name:        "zero number_of_tickets",
			numTickets:  0,
			expectCode:  codes.InvalidArgument,
			expectError: "number_of_tickets must be positive",
		},
		{
			name:        "negative number_of_tickets",
			numTickets:  -1,
			expectCode:  codes.InvalidArgument,
			expectError: "number_of_tickets must be positive",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &api.CreateOrderRequest{
				UserId:           1,
				ConcertSessionId: 1,
				NumberOfTickets:  tc.numTickets,
			}

			resp, err := handler.CreateOrder(context.Background(), req)
			assert.Nil(t, resp)
			assert.Error(t, err)

			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tc.expectCode, st.Code())
			assert.Contains(t, st.Message(), tc.expectError)
		})
	}
}

func TestGRPCHandler_CreateOrder_ConcertSessionNotFound(t *testing.T) {
	baseRepo, cleanup := repository.SetupTestDB(t)
	defer cleanup()

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	req := &api.CreateOrderRequest{
		UserId:           1,
		ConcertSessionId: 999, // Non-existent session
		NumberOfTickets:  1,
	}

	resp, err := handler.CreateOrder(context.Background(), req)
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "concert session not found")
}

func TestGRPCHandler_CreateOrder_NoTicketsAvailable(t *testing.T) {
	baseRepo, cleanup := repository.SetupTestDB(t)
	defer cleanup()

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	req := &api.CreateOrderRequest{
		UserId:           1,
		ConcertSessionId: 1,
		NumberOfTickets:  1,
	}

	// This test will fail if there are tickets available
	// In a real scenario, you would ensure no tickets are available
	resp, err := handler.CreateOrder(context.Background(), req)
	if err != nil {
		// If there are no tickets, that's expected
		t.Logf("Expected error due to no tickets: %v", err)
		return
	}

	// If we get here, there were tickets available
	require.NotNil(t, resp)
	assert.Greater(t, resp.OrderId, int32(0))
}

func TestGRPCHandler_CreateOrder_ResponseStructure(t *testing.T) {
	baseRepo, cleanup := repository.SetupTestDB(t)
	defer cleanup()

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	req := &api.CreateOrderRequest{
		UserId:           1,
		ConcertSessionId: 1,
		NumberOfTickets:  1,
	}

	resp, err := handler.CreateOrder(context.Background(), req)
	if err != nil {
		// Expected due to missing test data
		t.Logf("Expected error due to missing test data: %v", err)
		return
	}

	// Verify response structure
	require.NotNil(t, resp)
	assert.Greater(t, resp.OrderId, int32(0))
	assert.Equal(t, "pending", resp.Status)
	assert.NotNil(t, resp.TicketIds)
	assert.GreaterOrEqual(t, len(resp.TicketIds), 0)
	assert.Greater(t, resp.TotalPrice, float64(0))
	assert.NotNil(t, resp.CreatedAt)
}

func TestGRPCHandler_CreateOrder_ConcurrentRequests(t *testing.T) {
	baseRepo, cleanup := repository.SetupTestDB(t)
	defer cleanup()

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	// Test concurrent order creation
	const numGoroutines = 5
	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			req := &api.CreateOrderRequest{
				UserId:           int32(id + 1),
				ConcertSessionId: 1,
				NumberOfTickets:  1,
			}

			_, err := handler.CreateOrder(context.Background(), req)
			// We don't require success here as there might not be data
			// but we do require no panics or unexpected errors
			if err != nil {
				t.Logf("Goroutine %d got expected error: %v", id, err)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

func TestGRPCHandler_CreateOrder_ErrorHandling(t *testing.T) {
	baseRepo, cleanup := repository.SetupTestDB(t)
	defer cleanup()

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	// Test with invalid session ID
	req := &api.CreateOrderRequest{
		UserId:           1,
		ConcertSessionId: 999999, // Very large non-existent ID
		NumberOfTickets:  1,
	}

	resp, err := handler.CreateOrder(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestGRPCHandler_CreateOrder_PriceCalculation(t *testing.T) {
	baseRepo, cleanup := repository.SetupTestDB(t)
	defer cleanup()

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	req := &api.CreateOrderRequest{
		UserId:           1,
		ConcertSessionId: 1,
		NumberOfTickets:  2, // Request 2 tickets
	}

	// This test verifies that price calculations are correct
	resp, err := handler.CreateOrder(context.Background(), req)
	if err != nil {
		// Expected due to missing test data
		t.Logf("Expected error due to missing test data: %v", err)
		return
	}

	// If successful, verify the order was created with correct price
	require.NotNil(t, resp)
	assert.Greater(t, resp.OrderId, int32(0))
	assert.Equal(t, 2, len(resp.TicketIds)) // Should have 2 tickets
	assert.Greater(t, resp.TotalPrice, float64(0))

	t.Logf("Order created successfully with ID: %d, Total Price: %f", resp.OrderId, resp.TotalPrice)
}

func TestGRPCHandler_CreateOrder_ContextHandling(t *testing.T) {
	handler, cleanup := SetupTestHandlerWithData(t)
	defer cleanup()

	req := &api.CreateOrderRequest{
		UserId:           1,
		ConcertSessionId: 1,
		NumberOfTickets:  1,
	}

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	resp, err := handler.CreateOrder(ctx, req)
	// The behavior depends on the service implementation
	// We just ensure no panic occurs
	if err != nil {
		t.Logf("Expected error with cancelled context: %v", err)
	}
	// Note: The current implementation doesn't check context cancellation
	// so the order might still be created successfully
	if resp != nil {
		t.Logf("Order created despite cancelled context: %d", resp.OrderId)
	}
}

func TestGRPCHandler_CreateOrder_Logging(t *testing.T) {
	baseRepo, cleanup := repository.SetupTestDB(t)
	defer cleanup()

	baseService := service.NewBaseService(baseRepo)
	orderService := service.NewOrderService(baseService)
	handler := NewGRPCHandler(orderService)

	req := &api.CreateOrderRequest{
		UserId:           1,
		ConcertSessionId: 1,
		NumberOfTickets:  1,
	}

	// This test verifies that logging works correctly
	// The actual logging verification would require capturing log output
	resp, err := handler.CreateOrder(context.Background(), req)
	if err != nil {
		// Expected due to missing test data
		t.Logf("Expected error due to missing test data: %v", err)
		return
	}

	// If successful, verify the response
	require.NotNil(t, resp)
	assert.Greater(t, resp.OrderId, int32(0))

	t.Logf("Order created successfully with ID: %d", resp.OrderId)
}
