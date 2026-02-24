package app

import (
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
	grpc_client "OrderService/pkg/client/grpc"

	pb "github.com/erdedan1/protocol/proto/order_service/gen"
	pbLogger "github.com/erdedan1/shared/interceptors/logger"
	"github.com/erdedan1/shared/interceptors/recovery"
	requestid "github.com/erdedan1/shared/interceptors/request_id"
	log "github.com/erdedan1/shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	cfg        *config.Config
	grpcServer *grpc.Server
	L          log.Logger
}

func New(cfg *config.Config) *App {
	logger, err := log.NewLogger(cfg.Infrastructure.Log_LVL)
	if err != nil {
		panic(err)
	}
	return &App{
		cfg: cfg,
		L:   logger,
	}
}

const layer = "App"

func (a *App) Start() error {
	const method = "Start"
	a.L.Info(layer, method, "starting service")
	driver, clients, err := spot_instrument_service.New(a.cfg)
	if err != nil {
		return err
	}

	redisClient := cache.NewRedisClient(a.cfg)
	repos := usecase.NewRepositories(orderRepo.NewInMemoryRepo(a.L), user.NewRepo(a.L), market.NewMarketsCache(redisClient, a.L))
	srvs := usecase.NewServices(orderSrv.New(repos, a.L, driver))

	go func() {
		if err := a.startGRPCServer(*srvs); err != nil {
			os.Exit(1)
		}
	}()

	a.L.Info(layer, method, "service work")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	a.L.Info(layer, method, "waiting for shutdown signal")
	<-quit
	a.L.Info(layer, method, "shutdown signal received")
	err = a.mustCloseConnectionWithGRPCClients(clients)
	if err != nil {
		a.L.Error(layer, method, "failed to close connection with grpc clients", err)
	}
	a.stopGRPCServer()
	a.L.Info(layer, method, "service stopped gracefully")
	return nil
}

func (a *App) startGRPCServer(usecase usecase.Services) error {
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
	grpcHandler := order_service.New(usecase)
	pb.RegisterOrderServiceServer(grpcServer, grpcHandler)

	lis, err := net.Listen("tcp", a.cfg.GRPCServer.Address)
	if err != nil {
		a.L.Error(layer, method, "failed to listen: %v", err)
		return err
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		a.L.Error(layer, method, "grpc serve error", err)
		return err
	}
	return nil
}

func (a *App) stopGRPCServer() {
	const method = "stopGRPCServer"
	a.grpcServer.GracefulStop()
	a.L.Info(layer, method, "grpc server stopped gracefully")
}

func (a *App) mustCloseConnectionWithGRPCClients(clients []grpc_client.IGRPCClient) error {
	const method = "mustCloseConnectionWithGRPCClients"
	for _, client := range clients {
		if err := client.Close(); err != nil {
			a.L.Error(layer, method, "failed to close client", err)
		}
	}

	return nil
}
