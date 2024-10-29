package service

import (
	"context"

	"github.com/yaninyzwitty/sqs-go/internal/model"
	"github.com/yaninyzwitty/sqs-go/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order model.Order) (*model.Order, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	return s.repo.CreateOrder(ctx, order)
}
