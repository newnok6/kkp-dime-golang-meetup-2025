# Stock Order REST API

A RESTful API for managing stock orders with SQLite persistence.

## Features

- Create stock orders (Market and Limit orders)
- List all orders
- Get order details by ID
- Cancel pending orders
- SQLite database for persistence
- Hexagonal architecture (Ports & Adapters)

## Architecture

```
backend/
├── cmd/
│   └── main.go              # Application entry point
├── domain/
│   └── stock_order.go       # Domain models
├── port/
│   ├── stock_order_service.go    # Service interface
│   └── stock_order_repository.go # Repository interface
├── service/
│   └── stock_order_service.go    # Business logic implementation
└── adaptor/
    ├── http_handler.go           # HTTP handlers
    └── sqlite_repository.go      # SQLite implementation
```

## Getting Started

### Prerequisites

- Go 1.25.0 or higher
- SQLite3

### Installation

1. Install dependencies:
```bash
go mod download
```

2. Run the server:
```bash
go run cmd/main.go
```

The server will start on port 8080 (or the port specified in the PORT environment variable).

## API Endpoints

### Health Check
```bash
GET /health
```

### Create Order
```bash
POST /api/orders
Content-Type: application/json

{
  "symbol": "AAPL",
  "order_type": "MARKET",
  "order_side": "BUY",
  "quantity": 100,
  "price": 150.50
}
```

**Order Types:**
- `MARKET` - Market order (price is optional)
- `LIMIT` - Limit order (price is required)

**Order Sides:**
- `BUY` - Buy order
- `SELL` - Sell order

### Get Order by ID
```bash
GET /api/orders/{id}
```

### List All Orders
```bash
GET /api/orders
```

### Cancel Order
```bash
POST /api/orders/{id}/cancel
```

## Example Usage

### Create a Market Buy Order
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "AAPL",
    "order_type": "MARKET",
    "order_side": "BUY",
    "quantity": 100
  }'
```

### Create a Limit Sell Order
```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "GOOGL",
    "order_type": "LIMIT",
    "order_side": "SELL",
    "quantity": 50,
    "price": 2500.00
  }'
```

### Get All Orders
```bash
curl http://localhost:8080/api/orders
```

### Get Specific Order
```bash
curl http://localhost:8080/api/orders/{order-id}
```

### Cancel an Order
```bash
curl -X POST http://localhost:8080/api/orders/{order-id}/cancel
```

## Order Status

Orders go through the following statuses:

- `PENDING` - Order has been received and is being processed
- `FILLED` - Order has been executed (simulated after 2 seconds)
- `CANCELLED` - Order has been cancelled by the user
- `REJECTED` - Order was rejected (reserved for future use)

## Database

The application uses SQLite with the following schema:

```sql
CREATE TABLE stock_orders (
    id TEXT PRIMARY KEY,
    symbol TEXT NOT NULL,
    order_type TEXT NOT NULL,
    order_side TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    price REAL,
    status TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    description TEXT
);
```

The database file `stock_orders.db` will be created in the working directory when you run the application.

## Configuration

- **PORT**: Set the `PORT` environment variable to change the server port (default: 8080)

```bash
PORT=3000 go run cmd/main.go
```

## Development

### Build
```bash
go build -o stock-api cmd/main.go
```

### Run
```bash
./stock-api
```

## License

MIT
