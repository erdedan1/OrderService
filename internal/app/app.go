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
	"go.opentelemetry.io/otel"

	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	grpcServer *grpc.Server
	log        log.Logger
}

func New(cfg *config.Config, log log.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

const layer = "App"

func (a *App) Start(ctx context.Context) error {
	const method = "Start"
	tracer := otel.Tracer("order-service")
	_, span := tracer.Start(ctx, "service.start")
	span.End()
	a.log.Info(layer, method, "starting service")

	marketService, err := spot_instrument_service.NewMarketService(a.cfg)
	if err != nil {
		return err
	}
	defer marketService.Close()

	redisClient := cache.NewRedisClient(a.cfg)
	orderRepository := orderRepo.NewInMemoryRepo(a.log)
	orderStatusSubscriber := orderStatusRepo.NewRedisSubscriber(redisClient, a.log)
	orderStatusPublisher := orderStatusRepo.NewRedisPublisher(redisClient, a.log)
	userRepository := user.NewRepo(a.log)
	marketCache := market.NewMarketsCache(redisClient, a.log)
	orderService := orderSrv.New(
		orderRepository,
		userRepository,
		marketCache,
		marketService,
		orderStatusSubscriber,
		orderStatusPublisher,
		a.log,
	)

	serverErrCh := make(chan error, 1)

	grpcServer, err := order_service.NewGRPCServer(a.cfg.GRPCServer.Address, orderService, a.log)
	if err != nil {
		return err
	}

	go func() {
		serverErrCh <- grpcServer.Start()
	}()

	a.log.Info(layer, method, "service work")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	a.log.Info(layer, method, "waiting for shutdown signal")

	select {
	case <-ctx.Done():
		a.log.Info(layer, method, "context cancelled")
	case sig := <-quit:
		a.log.Info(layer, method, "shutdown signal received", "signal", sig.String())
	case err := <-serverErrCh:
		if err != nil {
			return err
		}
		return nil
	}

	grpcServer.Stop()

	if err := <-serverErrCh; err != nil && !order_service.IsExpectedStop(err) {
		return err
	}

	a.log.Info(layer, method, "service stopped gracefully")
	return nil
}
