package adaptor

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/domain"
	"github.com/newnok6/kkp-dime-golang-meetup-2025/backend/port"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (port.StockOrderRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &sqliteRepository{db: db}
	if err := repo.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

func (r *sqliteRepository) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS stock_orders (
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

	CREATE INDEX IF NOT EXISTS idx_symbol ON stock_orders(symbol);
	CREATE INDEX IF NOT EXISTS idx_status ON stock_orders(status);
	CREATE INDEX IF NOT EXISTS idx_created_at ON stock_orders(created_at DESC);
	`

	_, err := r.db.Exec(query)
	return err
}

func (r *sqliteRepository) Create(ctx context.Context, order *domain.StockOrder) error {
	query := `
		INSERT INTO stock_orders (id, symbol, order_type, order_side, quantity, price, status, created_at, updated_at, description)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		order.ID,
		order.Symbol,
		order.OrderType,
		order.OrderSide,
		order.Quantity,
		order.Price,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
		order.Description,
	)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

func (r *sqliteRepository) GetByID(ctx context.Context, orderID string) (*domain.StockOrder, error) {
	query := `
		SELECT id, symbol, order_type, order_side, quantity, price, status, created_at, updated_at, description
		FROM stock_orders
		WHERE id = ?
	`

	order := &domain.StockOrder{}
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.ID,
		&order.Symbol,
		&order.OrderType,
		&order.OrderSide,
		&order.Quantity,
		&order.Price,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.Description,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

func (r *sqliteRepository) List(ctx context.Context) ([]*domain.StockOrder, error) {
	query := `
		SELECT id, symbol, order_type, order_side, quantity, price, status, created_at, updated_at, description
		FROM stock_orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	defer rows.Close()

	orders := []*domain.StockOrder{}
	for rows.Next() {
		order := &domain.StockOrder{}
		err := rows.Scan(
			&order.ID,
			&order.Symbol,
			&order.OrderType,
			&order.OrderSide,
			&order.Quantity,
			&order.Price,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

func (r *sqliteRepository) Update(ctx context.Context, order *domain.StockOrder) error {
	query := `
		UPDATE stock_orders
		SET symbol = ?, order_type = ?, order_side = ?, quantity = ?, price = ?,
		    status = ?, updated_at = ?, description = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query,
		order.Symbol,
		order.OrderType,
		order.OrderSide,
		order.Quantity,
		order.Price,
		order.Status,
		order.UpdatedAt,
		order.Description,
		order.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found: %s", order.ID)
	}

	return nil
}

func (r *sqliteRepository) Close() error {
	return r.db.Close()
}
