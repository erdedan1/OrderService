package app

import (
	"context"
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

func Build(cfg *config.Config, log log.Logger) (*App, error) {

	redis := cache.NewRedisClient(cfg)

	orderRepo := orderRepo.NewInMemoryRepo(log)
	userRepo := user.NewRepo(log)

	subscriber := orderStatusRepo.NewRedisSubscriber(redis, log)
	publisher := orderStatusRepo.NewRedisPublisher(redis, log)

	marketService, err := spot_instrument_service.NewMarketService(cfg)
	if err != nil {
		return nil, err
	}

	marketCache := market.NewMarketsCache(redis, log)

	orderService := orderSrv.New(
		orderRepo,
		userRepo,
		marketCache,
		marketService,
		subscriber,
		publisher,
		log,
	)

	grpcServer, err := order_service.NewGRPCServer("asd", orderService, log)
	if err != nil {
		return nil, err
	}

	return New(cfg, grpcServer, log), nil
}

func (a *App) Start(ctx context.Context) error {

	a.grpcServer.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case <-quit:
	}

	a.grpcServer.Stop()
	return nil
}
