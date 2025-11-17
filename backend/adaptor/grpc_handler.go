package adaptor

import (
	"context"

	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/domain"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/port"
	pb "github.com/newnok6/kkp-dime-golang-meetup-2025/backend/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	pb.UnimplementedStockOrderServiceServer
	service port.StockOrderService
}

func NewGRPCHandler(service port.StockOrderService) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

// CreateOrder handles the gRPC CreateOrder request
func (h *GRPCHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.StockOrder, error) {
	// Convert protobuf request to domain request
	domainReq := domain.CreateOrderRequest{
		Symbol:    req.Symbol,
		OrderType: convertProtoOrderTypeToDomain(req.OrderType),
		OrderSide: convertProtoOrderSideToDomain(req.OrderSide),
		Quantity:  int(req.Quantity),
		Price:     req.Price,
	}

	// Call service
	order, err := h.service.CreateOrder(ctx, domainReq)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create order: %v", err)
	}

	// Convert domain order to protobuf
	return convertDomainOrderToProto(order), nil
}

// GetOrder handles the gRPC GetOrder request
func (h *GRPCHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.StockOrder, error) {
	order, err := h.service.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	return convertDomainOrderToProto(order), nil
}

// ListOrders handles the gRPC ListOrders request
func (h *GRPCHandler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := h.service.ListOrders(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	// Convert domain orders to protobuf
	protoOrders := make([]*pb.StockOrder, 0, len(orders))
	for _, order := range orders {
		protoOrders = append(protoOrders, convertDomainOrderToProto(order))
	}

	return &pb.ListOrdersResponse{
		Orders: protoOrders,
	}, nil
}

// CancelOrder handles the gRPC CancelOrder request
func (h *GRPCHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	err := h.service.CancelOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to cancel order: %v", err)
	}

	return &pb.CancelOrderResponse{
		Message: "order cancelled successfully",
	}, nil
}

// Helper functions to convert between protobuf and domain types

func convertDomainOrderToProto(order *domain.StockOrder) *pb.StockOrder {
	return &pb.StockOrder{
		Id:          order.ID,
		Symbol:      order.Symbol,
		OrderType:   convertDomainOrderTypeToProto(order.OrderType),
		OrderSide:   convertDomainOrderSideToProto(order.OrderSide),
		Quantity:    int32(order.Quantity),
		Price:       order.Price,
		Status:      convertDomainOrderStatusToProto(order.Status),
		CreatedAt:   timestamppb.New(order.CreatedAt),
		UpdatedAt:   timestamppb.New(order.UpdatedAt),
		Description: order.Description,
	}
}

func convertDomainOrderTypeToProto(orderType domain.OrderType) pb.OrderType {
	switch orderType {
	case domain.OrderTypeMarket:
		return pb.OrderType_MARKET
	case domain.OrderTypeLimit:
		return pb.OrderType_LIMIT
	default:
		return pb.OrderType_ORDER_TYPE_UNSPECIFIED
	}
}

func convertProtoOrderTypeToDomain(orderType pb.OrderType) domain.OrderType {
	switch orderType {
	case pb.OrderType_MARKET:
		return domain.OrderTypeMarket
	case pb.OrderType_LIMIT:
		return domain.OrderTypeLimit
	default:
		return ""
	}
}

func convertDomainOrderSideToProto(orderSide domain.OrderSide) pb.OrderSide {
	switch orderSide {
	case domain.OrderSideBuy:
		return pb.OrderSide_BUY
	case domain.OrderSideSell:
		return pb.OrderSide_SELL
	default:
		return pb.OrderSide_ORDER_SIDE_UNSPECIFIED
	}
}

func convertProtoOrderSideToDomain(orderSide pb.OrderSide) domain.OrderSide {
	switch orderSide {
	case pb.OrderSide_BUY:
		return domain.OrderSideBuy
	case pb.OrderSide_SELL:
		return domain.OrderSideSell
	default:
		return ""
	}
}

func convertDomainOrderStatusToProto(orderStatus domain.OrderStatus) pb.OrderStatus {
	switch orderStatus {
	case domain.OrderStatusPending:
		return pb.OrderStatus_PENDING
	case domain.OrderStatusFilled:
		return pb.OrderStatus_FILLED
	case domain.OrderStatusCancelled:
		return pb.OrderStatus_CANCELLED
	case domain.OrderStatusRejected:
		return pb.OrderStatus_REJECTED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}
