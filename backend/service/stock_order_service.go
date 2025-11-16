package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/domain"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/port"
)

type stockOrderService struct {
	repo port.StockOrderRepository
}

func NewStockOrderService(repo port.StockOrderRepository) port.StockOrderService {
	return &stockOrderService{
		repo: repo,
	}
}

func (s *stockOrderService) CreateOrder(ctx context.Context, req domain.CreateOrderRequest) (*domain.StockOrder, error) {
	// Validate order type and price
	if req.OrderType == domain.OrderTypeLimit && req.Price <= 0 {
		return nil, fmt.Errorf("limit order must have a price greater than 0")
	}

	order := &domain.StockOrder{
		ID:        uuid.New().String(),
		Symbol:    req.Symbol,
		OrderType: req.OrderType,
		OrderSide: req.OrderSide,
		Quantity:  req.Quantity,
		Price:     req.Price,
		Status:    domain.OrderStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	go s.processOrder(ctx, order.ID)
	return order, nil
}

func (s *stockOrderService) GetOrder(ctx context.Context, orderID string) (*domain.StockOrder, error) {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *stockOrderService) ListOrders(ctx context.Context) ([]*domain.StockOrder, error) {
	orders, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *stockOrderService) CancelOrder(ctx context.Context, orderID string) error {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.Status != domain.OrderStatusPending {
		return fmt.Errorf("cannot cancel order with status: %s", order.Status)
	}

	order.Status = domain.OrderStatusCancelled
	order.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, order); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

// processOrder simulates order processing and updates the database
func (s *stockOrderService) processOrder(ctx context.Context, orderID string) {
	log.Printf("[processOrder] Starting background processing for order: %s", orderID)

	// Simulate processing delay
	log.Printf("[processOrder] Order %s: Simulating 2-second processing delay", orderID)
	time.Sleep(2 * time.Second)

	log.Printf("[processOrder] Order %s: Processing delay complete, fetching order from database", orderID)
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		log.Printf("[processOrder] ERROR: Order %s: Failed to fetch from database: %v", orderID, err)
		return
	}

	if order.Status != domain.OrderStatusPending {
		log.Printf("[processOrder] Order %s: Skipping update, status is already %s", orderID, order.Status)
		return
	}

	log.Printf("[processOrder] Order %s: Updating status from PENDING to FILLED", orderID)
	order.Status = domain.OrderStatusFilled
	order.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, order); err != nil {
		log.Printf("[processOrder] ERROR: Order %s: Failed to update in database: %v", orderID, err)
		return
	}

	log.Printf("[processOrder] SUCCESS: Order %s: Successfully marked as FILLED", orderID)
}
