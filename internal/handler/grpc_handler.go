package handler

import (
	"context"
	"time"

	"tickets/api"
	"tickets/internal/logger"
	"tickets/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCHandler implements the TicketsService gRPC interface
type GRPCHandler struct {
	api.UnimplementedTicketsServiceServer
	orderService *service.OrderService
	// Add other services as needed
}

// NewGRPCHandler creates a new gRPC handler
func NewGRPCHandler(orderService *service.OrderService) *GRPCHandler {
	return &GRPCHandler{
		orderService: orderService,
	}
}

// CreateOrder implements the CreateOrder gRPC method
func (h *GRPCHandler) CreateOrder(ctx context.Context, req *api.CreateOrderRequest) (*api.CreateOrderResponse, error) {
	logger.WithFields(map[string]interface{}{
		"user_id":            req.UserId,
		"concert_session_id": req.ConcertSessionId,
		"number_of_tickets":  req.NumberOfTickets,
	}).Info("Creating order via gRPC")

	// Validate request
	if req.UserId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be positive")
	}
	if req.ConcertSessionId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "concert_session_id must be positive")
	}
	if req.NumberOfTickets <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "number_of_tickets must be positive")
	}
	if req.NumberOfTickets > 3 {
		return nil, status.Errorf(codes.InvalidArgument, "maximum 3 tickets allowed per order")
	}

	// Convert gRPC request to service request
	serviceReq := &service.CreateOrderRequest{
		UserID:           int(req.UserId),
		ConcertSessionID: int(req.ConcertSessionId),
		NumberOfTickets:  int(req.NumberOfTickets),
	}

	// Call service layer
	serviceResp, err := h.orderService.CreateOrder(serviceReq)
	if err != nil {
		logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":            req.UserId,
			"concert_session_id": req.ConcertSessionId,
		}).Error("Failed to create order")

		// Convert service errors to gRPC status codes
		switch err.Error() {
		case "concert session not found":
			return nil, status.Errorf(codes.NotFound, "concert session not found")
		case "no tickets available":
			return nil, status.Errorf(codes.ResourceExhausted, "no tickets available")
		case "request cannot be nil":
			return nil, status.Errorf(codes.InvalidArgument, "request cannot be nil")
		case "number of tickets must be greater than 0":
			return nil, status.Errorf(codes.InvalidArgument, "number_of_tickets must be positive")
		case "maximum 3 tickets allowed per order":
			return nil, status.Errorf(codes.InvalidArgument, "maximum 3 tickets allowed per order")
		default:
			return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
		}
	}

	// Convert service response to gRPC response
	resp := &api.CreateOrderResponse{
		OrderId:    int32(serviceResp.OrderID),
		Status:     serviceResp.Status,
		TicketIds:  serviceResp.TicketIDs,
		TotalPrice: float64(serviceResp.TotalPrice.InexactFloat64()),
		CreatedAt:  timestamppb.New(time.Unix(serviceResp.CreatedAt/1000, 0)),
	}

	logger.WithFields(map[string]interface{}{
		"order_id": serviceResp.OrderID,
		"status":   serviceResp.Status,
		"tickets":  len(serviceResp.TicketIDs),
	}).Info("Order created successfully via gRPC")

	return resp, nil
}

// GetOrder implements the GetOrder gRPC method
// NOTE: This method is intentionally unimplemented as it's not part of the current scope
func (h *GRPCHandler) GetOrder(ctx context.Context, req *api.GetOrderRequest) (*api.GetOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "GetOrder not implemented")
}

// ListOrders implements the ListOrders gRPC method
// NOTE: This method is intentionally unimplemented as it's not part of the current scope
func (h *GRPCHandler) ListOrders(ctx context.Context, req *api.ListOrdersRequest) (*api.ListOrdersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ListOrders not implemented")
}

// GetConcertSession implements the GetConcertSession gRPC method
// NOTE: This method is intentionally unimplemented as it's not part of the current scope
func (h *GRPCHandler) GetConcertSession(ctx context.Context, req *api.GetConcertSessionRequest) (*api.GetConcertSessionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "GetConcertSession not implemented")
}

// ListConcertSessions implements the ListConcertSessions gRPC method
// NOTE: This method is intentionally unimplemented as it's not part of the current scope
func (h *GRPCHandler) ListConcertSessions(ctx context.Context, req *api.ListConcertSessionsRequest) (*api.ListConcertSessionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ListConcertSessions not implemented")
}

// GetAvailableTickets implements the GetAvailableTickets gRPC method
// NOTE: This method is intentionally unimplemented as it's not part of the current scope
func (h *GRPCHandler) GetAvailableTickets(ctx context.Context, req *api.GetAvailableTicketsRequest) (*api.GetAvailableTicketsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "GetAvailableTickets not implemented")
}
