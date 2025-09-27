package order_service

import (
	"context"
	"time"

	inmemory "github.com/deimossy/wb-order-service/internal/cache/in_memory"
	"github.com/deimossy/wb-order-service/internal/models"
	"github.com/deimossy/wb-order-service/internal/repository/postgres"
	customerrors "github.com/deimossy/wb-order-service/pkg/custom_errors"
	"go.uber.org/zap"
)

type OrderService interface {
	GetOrder(ctx context.Context, uid string) (*models.Order, error)
	SaveOrder(ctx context.Context, order *models.Order) error
	WarmUpCache(ctx context.Context, n int) error
}

type OrderHolder struct {
	Order *models.Order
}

type service struct {
	repo  postgres.OrderRepository
	cache *inmemory.MemoryCache[*OrderHolder]
	ttl   time.Duration
	logg   *zap.Logger
}

func NewService(repo postgres.OrderRepository, cache *inmemory.MemoryCache[*OrderHolder], logg *zap.Logger) OrderService {
	return &service{
		repo:  repo,
		cache: cache,
		ttl:   5 * time.Minute,
		logg:   logg,
	}
}

func (s *service) GetOrder(ctx context.Context, uid string) (*models.Order, error) {
	if h, ok := s.cache.Get(uid); ok {
		if h != nil && h.Order != nil {
			return h.Order, nil
		}
	}

	o, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, customerrors.ErrOrderNotFound
	}

	s.cache.Set(uid, &OrderHolder{Order: o}, s.ttl)

	return o, nil
}

func (s *service) SaveOrder(ctx context.Context, order *models.Order) error {
	if order == nil {
		return customerrors.ErrNilOrder
	}
	if err := s.repo.Save(ctx, order); err != nil {
		return err
	}

	s.cache.Set(order.OrderUID, &OrderHolder{Order: order}, s.ttl)

	return nil
}

func (s *service) WarmUpCache(ctx context.Context, n int) error {
	orders, err := s.repo.GetLastN(ctx, n)
	if err != nil {
		return err
	}
	for _, o := range orders {
		s.cache.Set(o.OrderUID, &OrderHolder{Order: o}, s.ttl)
	}

	return nil
}
