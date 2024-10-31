package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yaninyzwitty/sqs-go/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) (*model.Order, error)
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	query := `INSERT INTO orders(product_id, quantity, total_price) VALUES($1, $2, $3) RETURNING id, product_id, total_price, order_date`
	err := r.db.QueryRow(ctx, query, order.ProductId, order.Quantity, order.TotalPrice).Scan(&order.ID, &order.ProductId, &order.TotalPrice, &order.OrderDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return &order, nil

}
