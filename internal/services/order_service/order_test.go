package order_service_test

import (
	"context"
	"testing"
	"time"

	inmemory "github.com/deimossy/wb-order-service/internal/cache/in_memory"
	"github.com/deimossy/wb-order-service/internal/models"
	"github.com/deimossy/wb-order-service/internal/repository/postgres"
	"github.com/deimossy/wb-order-service/internal/services/order_service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockRepo struct {
	calledGetByID int
}

func (m *mockRepo) GetByID(ctx context.Context, uid string) (*models.Order, error) {
	m.calledGetByID++
	time.Sleep(1 * time.Millisecond)
	return &models.Order{OrderUID: uid}, nil
}

func (m *mockRepo) Save(ctx context.Context, order *models.Order) error {
	return nil
}

func (m *mockRepo) GetLastN(ctx context.Context, n int) ([]*models.Order, error) {
	return []*models.Order{{OrderUID: "last"}}, nil
}

func TestGetOrder_CacheHit(t *testing.T) {
	logg := zap.NewNop()
	cache := inmemory.NewMemoryCache[*order_service.OrderHolder]()
	var repo postgres.OrderRepository = &mockRepo{}
	svc := order_service.NewService(repo, cache, logg)

	ctx := context.Background()
	uid := "123"

	o1, err := svc.GetOrder(ctx, uid)
	assert.NoError(t, err)
	assert.Equal(t, uid, o1.OrderUID)
	assert.Equal(t, 1, repo.(*mockRepo).calledGetByID)

	o2, err := svc.GetOrder(ctx, uid)
	assert.NoError(t, err)
	assert.Equal(t, uid, o2.OrderUID)
	assert.Equal(t, 1, repo.(*mockRepo).calledGetByID, "репозиторий не должен вызываться второй раз")
}

func TestSaveOrder(t *testing.T) {
	logg := zap.NewNop()
	cache := inmemory.NewMemoryCache[*order_service.OrderHolder]()
	var repo postgres.OrderRepository = &mockRepo{}
	svc := order_service.NewService(repo, cache, logg)

	ctx := context.Background()
	order := &models.Order{OrderUID: "456"}

	err := svc.SaveOrder(ctx, order)
	assert.NoError(t, err)

	h, ok := cache.Get(order.OrderUID)
	assert.True(t, ok)
	assert.Equal(t, order, h.Order)
}

func TestWarmUpCache(t *testing.T) {
	logg := zap.NewNop()
	cache := inmemory.NewMemoryCache[*order_service.OrderHolder]()
	var repo postgres.OrderRepository = &mockRepo{}
	svc := order_service.NewService(repo, cache, logg)

	ctx := context.Background()
	err := svc.WarmUpCache(ctx, 10)
	assert.NoError(t, err)

	h, ok := cache.Get("last")
	assert.True(t, ok)
	assert.Equal(t, "last", h.Order.OrderUID)
}
