# KKP DIME Golang Meetup 2025

A demonstration project showcasing graceful shutdown patterns in Go, built for the KKP DIME Golang Meetup 2025. This project uses a Stock Order REST API as an example application to demonstrate different approaches to implementing graceful shutdown in production systems.

## Overview

This repository demonstrates three different implementations of graceful shutdown patterns in Go:

1. **Demo 1**: Basic graceful shutdown with concurrent cleanup
2. **Demo 2**: Extended graceful shutdown with pre-cleanup processing
3. **Demo 3**: Generic shutdown service with reusable shutdown logic

All demos use the same Stock Order REST API application, built with hexagonal architecture principles.

## Project Structure

```
kkp-dime-golang-meetup-2025/
├── backend/
│   ├── cmd/
│   │   ├── demo_1/         # Basic graceful shutdown
│   │   ├── demo_2/         # Shutdown with processing
│   │   └── demo_3/         # Generic shutdown pattern
│   ├── domain/             # Domain models
│   ├── port/               # Port interfaces
│   ├── service/            # Business logic
│   ├── adaptor/            # Adaptors (HTTP, SQLite)
│   ├── graceful-shutdown.d2
│   └── graceful-shutdown.svg
└── README.md
```

## Features

### Stock Order API
- Create stock orders (Market and Limit orders)
- List all orders
- Get order details by ID
- Cancel pending orders
- SQLite database for persistence
- Hexagonal architecture (Ports & Adapters)

### Graceful Shutdown Demonstrations

#### Demo 1: Basic Graceful Shutdown
Basic implementation using `sync.WaitGroup` to concurrently shut down the HTTP server and close the database connection.

**Key Features:**
- 30-second timeout context
- Concurrent shutdown of server and database
- Proper error handling for shutdown operations

#### Demo 2: Extended Graceful Shutdown
Enhanced version with additional processing logic before shutdown, simulating real-world cleanup tasks.

**Key Features:**
- Pre-shutdown processing tasks
- Concurrent execution of cleanup operations
- Demonstrates handling long-running shutdown tasks

#### Demo 3: Generic Shutdown Pattern
Reusable shutdown service that can handle any service implementing the `io.Closer` interface.

**Key Features:**
- Generic shutdown function
- Service registry pattern
- Type-safe shutdown handling with interface assertions
- Highly reusable across different projects

## Prerequisites

- Go 1.25.0 or higher
- SQLite3

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

## Running the Demos

### Demo 1: Basic Graceful Shutdown
```bash
cd backend
go run cmd/demo_1/main.go
```

### Demo 2: Extended Graceful Shutdown
```bash
cd backend
go run cmd/demo_2/main.go
```

### Demo 3: Generic Shutdown Pattern
```bash
cd backend
go run cmd/demo_3/main.go
```

Each demo starts the server on port 8082 (configurable via `PORT` environment variable).

To test graceful shutdown, send a SIGINT or SIGTERM signal:
```bash
# Press Ctrl+C in the terminal running the server
# OR
kill -SIGTERM <process-id>
```

## API Usage

For detailed API documentation, see [backend/README.md](backend/README.md).

### Quick Start

```bash
# Health check
curl http://localhost:8082/health

# Create a market buy order
curl -X POST http://localhost:8082/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "AAPL",
    "order_type": "MARKET",
    "order_side": "BUY",
    "quantity": 100
  }'

# List all orders
curl http://localhost:8082/api/orders
```

## Architecture

The application follows hexagonal architecture principles:

- **Domain Layer**: Core business logic and entities (`domain/`)
- **Port Layer**: Interfaces for services and repositories (`port/`)
- **Service Layer**: Business logic implementation (`service/`)
- **Adaptor Layer**: External integrations (HTTP handlers, SQLite repository) (`adaptor/`)

## Graceful Shutdown Flow

```
1. Receive SIGINT/SIGTERM signal
2. Stop accepting new requests
3. Execute cleanup tasks concurrently:
   - Shutdown HTTP server (wait for active connections)
   - Close database connections
   - (Demo 2/3) Execute additional cleanup logic
4. Wait for all cleanup tasks to complete (with timeout)
5. Exit application
```

## Technologies Used

- **Go 1.25.0**: Programming language
- **Gorilla Mux**: HTTP router
- **SQLite**: Database
- **UUID**: Unique identifier generation

## Learning Objectives

This project demonstrates:
- Implementing graceful shutdown in Go applications
- Using `context.Context` for timeout management
- Concurrent cleanup operations with `sync.WaitGroup`
- Signal handling in Go (`os/signal` package)
- Hexagonal architecture pattern
- REST API development in Go
- SQLite integration

## Configuration

- **PORT**: Server port (default: 8082)

Example:
```bash
PORT=3000 go run cmd/demo_1/main.go
```

## Development

### Build
```bash
cd backend
go build -o stock-api cmd/demo_1/main.go
```

### Run
```bash
./stock-api
```

## License

MIT
