package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deimossy/wb-order-service/internal/handlers"
	"github.com/deimossy/wb-order-service/internal/mapper"
	"github.com/deimossy/wb-order-service/internal/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockService struct{}

func (m *mockService) GetOrder(ctx context.Context, uid string) (*models.Order, error) {
	return &models.Order{OrderUID: uid}, nil
}
func (m *mockService) SaveOrder(ctx context.Context, order *models.Order) error { return nil }
func (m *mockService) WarmUpCache(ctx context.Context, n int) error             { return nil }

func TestGetOrderHandler(t *testing.T) {
	svc := &mockService{}
	logger := zap.NewNop()

	h := handlers.NewOrderHandler(svc, logger)

	req := httptest.NewRequest(http.MethodGet, "/order/123", nil)
	req = mux.SetURLVars(req, map[string]string{"order_uid": "123"})

	w := httptest.NewRecorder()

	h.GetOrder(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var order mapper.OrderResp
	err := json.NewDecoder(resp.Body).Decode(&order)
	assert.NoError(t, err)
	assert.Equal(t, "123", order.OrderUID)
}
