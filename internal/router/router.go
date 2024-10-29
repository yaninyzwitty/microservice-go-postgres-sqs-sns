package router

import (
	"net/http"

	"github.com/yaninyzwitty/sqs-go/internal/controller"
)

func NewRouter(controller controller.OrderController) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /orders", controller.CreateOrderHandler)
	return router
}
