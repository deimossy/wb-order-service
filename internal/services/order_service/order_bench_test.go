package order_service_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	inmemory "github.com/deimossy/wb-order-service/internal/cache/in_memory"
	"github.com/deimossy/wb-order-service/internal/models"
	"github.com/deimossy/wb-order-service/internal/services/order_service"
	"go.uber.org/zap"
)

func LoadOrderFromJSON(filename string) (*models.Order, error) {
	data, err := os.ReadFile("testdata/" + filename)
	if err != nil {
		return nil, err
	}

	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func BenchmarkOrderService_WithCache(b *testing.B) {
	logg := zap.NewNop()
	cache := inmemory.NewMemoryCache[*order_service.OrderHolder]()
	repo := &mockRepo{}
	svc := order_service.NewService(repo, cache, logg)

	ctx := context.Background()

	order, err := LoadOrderFromJSON("order.json")
	if err != nil {
		b.Fatalf("failed to load order: %v", err)
	}

	for i := 0; i < 30; i++ {
		order.OrderUID = fmt.Sprintf("%s_%d", order.OrderUID, i)
		_ = svc.SaveOrder(ctx, order)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.GetOrder(ctx, order.OrderUID+"_15")
	}
}

func BenchmarkOrderService_NoCache(b *testing.B) {
	cache := inmemory.NewMemoryCache[*order_service.OrderHolder]()
	repo := &mockRepo{}
	logg := zap.NewNop()
	_ = order_service.NewService(repo, cache, logg)

	ctx := context.Background()
	order, err := LoadOrderFromJSON("order.json")
	if err != nil {
		b.Fatalf("failed to load order: %v", err)
	}

	for i := 0; i < 30; i++ {
		order.OrderUID = fmt.Sprintf("%s_%d", order.OrderUID, i)
		_, _ = repo.GetByID(ctx, order.OrderUID)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetByID(ctx, order.OrderUID+"_15")
	}
}
