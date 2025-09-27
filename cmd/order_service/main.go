package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	inmemory "github.com/deimossy/wb-order-service/internal/cache/in_memory"
	"github.com/deimossy/wb-order-service/internal/config"
	"github.com/deimossy/wb-order-service/internal/handlers"
	"github.com/deimossy/wb-order-service/internal/kafka/consumer"
	"github.com/deimossy/wb-order-service/internal/middleware"
	"github.com/deimossy/wb-order-service/internal/repository/postgres"
	"github.com/deimossy/wb-order-service/internal/services/order_service"
	"github.com/deimossy/wb-order-service/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logg := logger.New()
	defer logg.Sync()

	cfg := config.LoadConfig(logg)

	db, err := sqlx.Connect("postgres", cfg.PgDSN)
	if err != nil {
		logg.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer db.Close()

	repo := postgres.NewRepository(db)

	cache := inmemory.NewMemoryCache[*order_service.OrderHolder]()

	svc := order_service.NewService(repo, cache, logg)

	if err := svc.WarmUpCache(context.Background(), cfg.CacheWarmSize); err != nil {
		logg.Warn("cache warmup failed", zap.Error(err))
	}

	router := handlers.NewRouter(svc, logg)
	handler := middleware.Recovery(router, logg)

	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           handler,
		ReadTimeout:       cfg.ServerReadTimeout,
		WriteTimeout:      cfg.ServerWriteTimeout,
		IdleTimeout:       cfg.ServerIdleTimeout,
		ReadHeaderTimeout: cfg.ServerReadHeaderTimeout,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logg.Info("server starting", zap.String("port", cfg.ServerPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatal("server failed", zap.Error(err))
		}
	}()

	go func() {
		c := consumer.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopic, "order_service_group", svc, logg)
		c.Run(ctx)
	}()

	<-stop
	logg.Info("shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logg.Error("server shutdown failed", zap.Error(err))
	}

	cancel()

	logg.Info("server gracefully stopped")
}
