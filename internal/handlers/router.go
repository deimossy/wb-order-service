package handlers

import (
	"net/http"

	"github.com/deimossy/wb-order-service/internal/services/order_service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(svc order_service.OrderService, logger *zap.Logger) http.Handler {
	r := mux.NewRouter()

	h := NewOrderHandler(svc, logger)

	r.HandleFunc("/order/{order_uid}", h.GetOrder).Methods(http.MethodGet)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/")))

	return r
}
