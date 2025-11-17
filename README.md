# KKP DIME Golang Meetup 2025

A comprehensive demonstration project showcasing graceful shutdown patterns in Go, built for the KKP DIME Golang Meetup 2025. This project implements a Stock Order system with both REST and gRPC APIs to demonstrate different approaches to graceful shutdown in production-grade distributed systems.

## Overview

This repository demonstrates three progressive implementations of graceful shutdown patterns in Go:

1. **Demo 1**: Basic graceful shutdown with HTTP REST API
2. **Demo 2**: Fail graceful shutdown usecase seperate thead
3. **Demo 3**: Advanced graceful shutdown managing both HTTP and gRPC servers and other dependencies

All demos implement the same Stock Order business logic, built using hexagonal architecture principles with dual API interfaces (REST and gRPC).

## Project Structure

```
kkp-dime-golang-meetup-2025/
├── backend/
│   ├── cmd/
│   │   ├── demo_1/         # Basic graceful shutdown (HTTP only)
│   │   ├── demo_2/         # Shutdown with processing (HTTP only)
│   │   ├── demo_3/         # Generic shutdown pattern (HTTP + gRPC)
│   │   ├── client/         # HTTP REST API client
│   │   └── client_grpc/    # gRPC client
│   ├── domain/             # Domain models and business entities
│   ├── port/               # Port interfaces (dependency inversion)
│   ├── service/            # Business logic implementation
│   ├── adaptor/            # Adaptors (HTTP, gRPC, SQLite)
│   ├── proto/              # Protocol Buffers definitions
│   └── stock_orders.db     # SQLite database
└── README.md
```

## Features

### Stock Order System
- **Order Management**: Create, retrieve, list, and cancel stock orders
- **Order Types**: Market and Limit orders
- **Order Sides**: Buy and Sell operations
- **Order Status Tracking**: Pending, Filled, Cancelled, Rejected
- **Dual API Support**: Both REST (HTTP/JSON) and gRPC (Protocol Buffers)
- **Persistent Storage**: SQLite database with automatic schema migration
- **Clean Architecture**: Hexagonal architecture (Ports & Adapters pattern)

### Graceful Shutdown Demonstrations

#### Demo 1: Basic Graceful Shutdown
Foundation implementation demonstrating concurrent cleanup of HTTP server and database connections using `sync.WaitGroup`.

**Key Features:**
- 30-second timeout context for shutdown operations
- Concurrent shutdown of HTTP server and database
- Signal handling (SIGINT/SIGTERM)
- Proper error handling and logging

#### Demo 2: Extended Graceful Shutdown
Enhanced version with additional processing logic before shutdown, simulating real-world scenarios like completing in-flight transactions or flushing caches. So that, It make database closed before HTTP Server shurdown

**Key Features:**
- Pre-shutdown processing hooks
- Concurrent execution of multiple cleanup operations
- Demonstrates handling of long-running shutdown tasks
- Graceful degradation patterns

#### Demo 3: Generic Shutdown Pattern with Multi-Protocol Support
Production-grade implementation featuring a reusable shutdown service that manages both HTTP and gRPC servers concurrently.

**Key Features:**
- Generic shutdown function using `io.Closer` interface
- Service registry pattern for managing multiple services
- Dual protocol support (HTTP REST + gRPC)
- Type-safe shutdown handling with interface assertions
- Highly reusable across different projects and service types
- Demonstrates shutdown coordination for heterogeneous service architecture

## Prerequisites

- Go 1.25.0 or higher
- SQLite3
- Protocol Buffers compiler (protoc) - for regenerating .proto files
- gRPC tools for Go

## Installation

1. Clone the repository:
```bash
git clone https://github.com/newnok6/kkp-dime-golang-meetup-2025.git
cd kkp-dime-golang-meetup-2025
```

2. Install dependencies:
```bash
cd backend
go mod download
```

3. (Optional) Regenerate Protocol Buffers:
```bash
cd backend
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/stock_order.proto
```

## Running the Demos

### Demo 1: Basic Graceful Shutdown (HTTP Only)
```bash
cd backend
go run cmd/demo_1/main.go
```

### Demo 2: Extended Graceful Shutdown (HTTP Only)
```bash
cd backend
go run cmd/demo_2/main.go
```

### Demo 3: Generic Shutdown Pattern (HTTP + gRPC)
```bash
cd backend
go run cmd/demo_3/main.go
```

**Demo 3 runs both servers:**
- HTTP server on port 8082 (configurable via `PORT` environment variable)
- gRPC server on port 50051 (configurable via `GRPC_PORT` environment variable)

### Testing Graceful Shutdown

Send a SIGINT or SIGTERM signal to trigger graceful shutdown:
```bash
# Press Ctrl+C in the terminal running the server
# OR
kill -SIGTERM <process-id>
```

Observe the logs to see the graceful shutdown sequence in action.

## API Usage

### REST API Examples

#### Health Check
```bash
curl http://localhost:8082/health
```

#### Create a Market Buy Order
```bash
curl -X POST http://localhost:8082/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "AAPL",
    "order_type": "MARKET",
    "order_side": "BUY",
    "quantity": 100
  }'
```

#### Create a Limit Sell Order
```bash
curl -X POST http://localhost:8082/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "GOOGL",
    "order_type": "LIMIT",
    "order_side": "SELL",
    "quantity": 50,
    "price": 150.75
  }'
```

#### List All Orders
```bash
curl http://localhost:8082/api/orders
```

#### Get Order by ID
```bash
curl http://localhost:8082/api/orders/{order-id}
```

#### Cancel Order
```bash
curl -X POST http://localhost:8082/api/orders/{order-id}/cancel
```

### gRPC API Examples

Use the provided gRPC client or tools like `grpcurl`:

#### Using the Demo Client
```bash
cd backend
go run cmd/client_grpc/main.go
```

#### Using grpcurl

List services:
```bash
grpcurl -plaintext localhost:50051 list
```

Create an order:
```bash
grpcurl -plaintext -d '{
  "symbol": "AAPL",
  "order_type": "MARKET",
  "order_side": "BUY",
  "quantity": 100
}' localhost:50051 stockorder.StockOrderService/CreateOrder
```

List orders:
```bash
grpcurl -plaintext -d '{}' localhost:50051 stockorder.StockOrderService/ListOrders
```

## Architecture

The application follows hexagonal architecture (Ports and Adapters) principles:

### Core Layers

- **Domain Layer** (`domain/`): Core business logic and entities
  - `StockOrder`: Main business entity
  - `OrderType`, `OrderSide`, `OrderStatus`: Value objects
  - Business rules and validations

- **Port Layer** (`port/`): Interfaces defining contracts
  - `StockOrderService`: Business logic interface
  - `StockOrderRepository`: Data persistence interface
  - Enables dependency inversion and testability

- **Service Layer** (`service/`): Business logic implementation
  - Order creation, validation, and processing
  - Implements port interfaces
  - Technology-agnostic business rules

- **Adaptor Layer** (`adaptor/`): External integrations
  - **HTTP Handler**: REST API implementation with Gorilla Mux
  - **gRPC Handler**: gRPC service implementation
  - **SQLite Repository**: Data persistence with SQLite
  - **Protocol Buffers**: Service definitions in `proto/`

### Architecture Benefits

- **Testability**: Mock ports for unit testing
- **Flexibility**: Easily swap implementations (e.g., PostgreSQL instead of SQLite)
- **Maintainability**: Clear separation of concerns
- **Scalability**: Add new adaptors (e.g., GraphQL) without changing core logic
- **Multi-Protocol Support**: Same business logic serves both REST and gRPC

## Graceful Shutdown Flow

```
1. Receive SIGINT/SIGTERM signal
2. Stop accepting new requests (both HTTP and gRPC)
3. Execute cleanup tasks concurrently:
   - Shutdown HTTP server (wait for active connections to complete)
   - Shutdown gRPC server (graceful stop with deadline)
   - Close database connections
   - (Demo 2/3) Execute additional cleanup logic
4. Wait for all cleanup tasks to complete (with 30s timeout)
5. Log shutdown completion
6. Exit application
```

### Shutdown Coordination

Demo 3 demonstrates coordinating shutdown across multiple server types:
- HTTP server: Uses `Shutdown()` with context timeout
- gRPC server: Uses `GracefulStop()` with goroutine timeout
- Database: Closes connections after servers stop accepting requests

## Technologies Used

- **Go 1.25.0**: Programming language
- **Gorilla Mux**: HTTP router and middleware
- **gRPC**: High-performance RPC framework
- **Protocol Buffers**: Interface Definition Language (IDL)
- **SQLite**: Embedded database
- **UUID**: Unique identifier generation (google/uuid)

## Configuration

### Environment Variables

- `PORT`: HTTP server port (default: 8082)
- `GRPC_PORT`: gRPC server port (default: 50051)

### Examples

Run with custom ports:
```bash
PORT=3000 GRPC_PORT=9090 go run cmd/demo_3/main.go
```

Run HTTP only (Demo 1 or 2):
```bash
PORT=8080 go run cmd/demo_1/main.go
```

## Learning Objectives

This project demonstrates best practices for:

### Graceful Shutdown Patterns
- Implementing graceful shutdown in Go applications
- Using `context.Context` for timeout management
- Concurrent cleanup operations with `sync.WaitGroup`
- Signal handling with `os/signal` package
- Managing multiple server lifecycles simultaneously

### Software Architecture
- Hexagonal architecture (Ports & Adapters pattern)
- Dependency inversion principle
- Clean separation of concerns
- Multi-protocol API design

### Go Development
- REST API development with Gorilla Mux
- gRPC service implementation
- Protocol Buffers usage
- SQLite integration
- Middleware patterns
- Error handling and logging

### Production Readiness
- Graceful shutdown for zero-downtime deployments
- Resource cleanup and connection management
- Timeout handling and deadline propagation
- Multi-service orchestration

## Development

### Build Demos

```bash
cd backend

# Build demo 1
go build -o demo1 cmd/demo_1/main.go

# Build demo 2
go build -o demo2 cmd/demo_2/main.go

# Build demo 3
go build -o demo3 cmd/demo_3/main.go
```

### Run Compiled Binary

```bash
./demo1
./demo2
./demo3
```

### Run Tests

```bash
cd backend
go test ./...
```

### Generate Protocol Buffers

After modifying `.proto` files:
```bash
cd backend
protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/stock_order.proto
```

## Project Highlights

- **Progressive Complexity**: Learn graceful shutdown from basic to advanced patterns
- **Real-World Application**: Stock order system resembles actual trading systems
- **Multi-Protocol**: Demonstrates modern microservice communication patterns
- **Production Patterns**: Code ready for production use
- **Educational**: Well-commented code with clear separation of concerns
- **Extensible**: Easy to add new features, protocols, or storage backends

## Troubleshooting

### Port Already in Use

If you see "address already in use" errors:
```bash
# Find process using port 8082
lsof -i :8082

# Find process using port 50051
lsof -i :50051

# Kill the process
kill -9 <PID>
```

### Database Locked

If SQLite database is locked:
```bash
# Remove the database file and restart
rm backend/stock_orders.db
```

### gRPC Connection Refused

Ensure Demo 3 is running (demos 1 and 2 don't start gRPC server):
```bash
go run cmd/demo_3/main.go
```

## Contributing

This is a demonstration project for educational purposes. Feel free to fork and experiment with different graceful shutdown patterns.

## License

MIT

## Resources

- [Go Context Package](https://pkg.go.dev/context)
- [gRPC Go Quick Start](https://grpc.io/docs/languages/go/quickstart/)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Graceful Shutdown in Go](https://pkg.go.dev/net/http#Server.Shutdown)

---

Built with care for KKP DIME Golang Meetup 2025
