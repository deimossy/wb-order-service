package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/deimossy/wb-order-service/internal/mapper"
	"github.com/deimossy/wb-order-service/internal/services/order_service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type OrderHandler struct {
	svc    order_service.OrderService
	logger *zap.Logger
}

func NewOrderHandler(svc order_service.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{svc: svc, logger: logger}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["order_uid"]

	ctx := r.Context()
	order, err := h.svc.GetOrder(ctx, uid)
	if err != nil {
		h.logger.Warn("order not found", zap.String("order_uid", uid), zap.Error(err))
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	resp := mapper.ModelToResponse(order)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	h.logger.Info("order retrieved", zap.String("order_uid", uid))
}
