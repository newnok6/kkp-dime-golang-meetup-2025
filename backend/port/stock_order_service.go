package port

import (
	"context"

	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/domain"
)

type StockOrderService interface {
	CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (*domain.StockOrder, error)
	GetOrder(ctx context.Context, orderID string) (*domain.StockOrder, error)
	ListOrders(ctx context.Context) ([]*domain.StockOrder, error)
	CancelOrder(ctx context.Context, orderID string) error
}
