package app

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"OrderService/config"
	"OrderService/internal/grpc/order_service"
	"OrderService/internal/grpc/spot_instrument_service"
	"OrderService/internal/repository/market"
	orderRepo "OrderService/internal/repository/order"
	"OrderService/internal/repository/user"
	orderSrv "OrderService/internal/service/order"
	"OrderService/internal/usecase"
	"OrderService/pkg/cache"

	pbOrder "github.com/erdedan1/protocol/proto/order_service/gen"
	pbLogger "github.com/erdedan1/shared/interceptors/logger"
	"github.com/erdedan1/shared/interceptors/recovery"
	requestid "github.com/erdedan1/shared/interceptors/request_id"
	log "github.com/erdedan1/shared/logger"
	"go.opentelemetry.io/otel"

	"go.uber.org/zap"
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
	userRepository := user.NewRepo(a.log)
	marketCache := market.NewMarketsCache(redisClient, a.log)
	orderService := orderSrv.New(
		orderRepository,
		userRepository,
		marketCache,
		marketService,
		a.log,
	)

	serverErrCh := make(chan error, 1)

	go func() {
		serverErrCh <- a.startGRPCServer(orderService)
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

	a.stopGRPCServer()
	if err := <-serverErrCh; err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}

	a.log.Info(layer, method, "service stopped gracefully")
	return nil
}

func (a *App) startGRPCServer(orderService usecase.OrderService) error {
	const method = "startGRPCServer"
	zap, _ := zap.NewProduction()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestid.XRequestIDServerInterceptor(),
			pbLogger.LoggerServerInterceptor(zap),
			recovery.RecoveryServerInterceptor(zap),
		),
	)

	a.grpcServer = grpcServer
	grpcHandler := order_service.New(orderService, a.log)
	pbOrder.RegisterOrderServiceServer(grpcServer, grpcHandler)

	lis, err := net.Listen("tcp", a.cfg.GRPCServer.Address)
	if err != nil {
		a.log.Error(layer, method, "failed to listen", err)
		return err
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		a.log.Error(layer, method, "grpc serve error", err)
		return err
	}
	return nil
}

func (a *App) stopGRPCServer() {
	const method = "stopGRPCServer"
	if a.grpcServer == nil {
		return
	}
	a.grpcServer.GracefulStop()
	a.log.Info(layer, method, "grpc server stopped gracefully")
}
