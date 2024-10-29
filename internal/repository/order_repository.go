package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yaninyzwitty/sqs-go/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) (*model.Order, error)
}

type orderRepository struct {
	snsClient   *sns.Client
	snsTopicArn *string
	db          *pgxpool.Pool
}

func NewOrderRepository(snsClient *sns.Client, topicArn *string, db *pgxpool.Pool) OrderRepository {
	return &orderRepository{snsClient: snsClient, snsTopicArn: topicArn, db: db}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order model.Order) (*model.Order, error) {
	return nil, nil
}
