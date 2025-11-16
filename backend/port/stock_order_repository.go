package port

import (
	"context"

	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/domain"
)

type StockOrderRepository interface {
	Create(ctx context.Context, order *domain.StockOrder) error
	GetByID(ctx context.Context, orderID string) (*domain.StockOrder, error)
	List(ctx context.Context) ([]*domain.StockOrder, error)
	Update(ctx context.Context, order *domain.StockOrder) error
	Close() error
}
