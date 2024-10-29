package controller

import (
	"encoding/json"
	"net/http"

	"github.com/yaninyzwitty/sqs-go/internal/model"
	"github.com/yaninyzwitty/sqs-go/internal/service"
)

type OrderController struct {
	service service.OrderService
}

func NewOrderController(service service.OrderService) *OrderController {
	return &OrderController{service: service}
}

func (c *OrderController) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var order model.Order
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if order.Quantity == 0 || order.Status == "" {
		http.Error(w, "Both quantity and status are required", http.StatusBadRequest)
		return
	}

	createdOrder, err := c.service.CreateOrder(ctx, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(createdOrder)
	if err != nil {
		http.Error(w, "Error marshalling to json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
