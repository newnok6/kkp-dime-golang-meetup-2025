package domain

import "time"

type OrderType string
type OrderSide string
type OrderStatus string

const (
	OrderTypeMarket OrderType = "MARKET"
	OrderTypeLimit  OrderType = "LIMIT"

	OrderSideBuy  OrderSide = "BUY"
	OrderSideSell OrderSide = "SELL"

	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusFilled    OrderStatus = "FILLED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
	OrderStatusRejected  OrderStatus = "REJECTED"
)

type StockOrder struct {
	ID          string      `json:"id"`
	Symbol      string      `json:"symbol"`
	OrderType   OrderType   `json:"order_type"`
	OrderSide   OrderSide   `json:"order_side"`
	Quantity    int         `json:"quantity"`
	Price       float64     `json:"price,omitempty"`
	Status      OrderStatus `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Description string      `json:"description,omitempty"`
}

type CreateOrderRequest struct {
	Symbol    string    `json:"symbol" validate:"required"`
	OrderType OrderType `json:"order_type" validate:"required"`
	OrderSide OrderSide `json:"order_side" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
	Price     float64   `json:"price,omitempty"`
}
