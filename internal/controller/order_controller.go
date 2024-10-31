package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
	"github.com/yaninyzwitty/sqs-go/internal/model"
	"github.com/yaninyzwitty/sqs-go/internal/service"
)

type OrderController struct {
	service   service.OrderService
	sqsClient *sqs.Client
	queueURL  *string
}

func NewOrderController(service service.OrderService, client *sqs.Client, queueURL *string) *OrderController {
	return &OrderController{service: service, sqsClient: client, queueURL: queueURL}
}

func (c *OrderController) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var order model.Order
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if order.Quantity == 0 {
		http.Error(w, " quantity is required", http.StatusBadRequest)
		return
	}
	order.ID = uuid.New()

	createdOrder, err := c.service.CreateOrder(ctx, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(createdOrder)
	if err != nil {
		http.Error(w, "Failed to marshal order response to JSON", http.StatusInternalServerError)
		return
	}

	_, err = c.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    c.queueURL,
		MessageBody: aws.String(string(response)),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send message to queue: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
