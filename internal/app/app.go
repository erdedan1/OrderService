package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"OrderService/config"
	"OrderService/internal/grpc/order_service"
	"OrderService/internal/grpc/spot_instrument_service"
	"OrderService/internal/repository/market"
	orderRepo "OrderService/internal/repository/order"
	orderStatusRepo "OrderService/internal/repository/order_status"
	"OrderService/internal/repository/user"
	orderSrv "OrderService/internal/service/order"
	"OrderService/pkg/cache"

	log "github.com/erdedan1/shared/logger"
	"go.opentelemetry.io/otel/sdk/trace"
)

type App struct {
	cfg        *config.Config
	grpcServer *order_service.GRPCServer
	log        log.Logger
}

func New(cfg *config.Config, grpcServer *order_service.GRPCServer, log log.Logger) *App {
	return &App{
		cfg:        cfg,
		grpcServer: grpcServer,
		log:        log,
	}
}

func Build(cfg *config.Config, log log.Logger, tp *trace.TracerProvider) (*App, error) {

	redis := cache.NewRedisClient(cfg)

	orderRepo := orderRepo.NewInMemoryRepo(log, tp)
	userRepo := user.NewRepo(log)

	subscriber := orderStatusRepo.NewRedisSubscriber(redis, log, tp)
	publisher := orderStatusRepo.NewRedisPublisher(redis, log, tp)

	marketService, err := spot_instrument_service.NewMarketService(cfg, tp)
	if err != nil {
		return nil, err
	}

	marketCache := market.NewMarketsCache(redis, log, tp)

	orderService := orderSrv.New(
		orderRepo,
		userRepo,
		marketCache,
		marketService,
		subscriber,
		publisher,
		log,
		tp,
	)

	grpcServer, err := order_service.NewGRPCServer(cfg.GRPCServer.Address, orderService, log)
	if err != nil {
		return nil, err
	}

	return New(cfg, grpcServer, log), nil
}

func (a *App) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.grpcServer.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("grpc server start failed: %w", err)
		}
		return nil
	case <-ctx.Done():
		a.grpcServer.Stop()
		return ctx.Err()
	case <-quit:
		a.grpcServer.Stop()
		if err := <-errCh; err != nil && !order_service.IsExpectedStop(err) {
			return fmt.Errorf("grpc server stop failed: %w", err)
		}
		return nil
	}
}
