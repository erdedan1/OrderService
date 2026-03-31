package order_service

import (
	"errors"
	"net"

	"OrderService/config"
	"OrderService/internal/usecase"

	pbOrder "github.com/erdedan1/protocol/proto/order_service/gen/v1"
	"github.com/erdedan1/shared/errs"
	pbLogger "github.com/erdedan1/shared/interceptors/logger"
	"github.com/erdedan1/shared/interceptors/recovery"
	requestid "github.com/erdedan1/shared/interceptors/request_id"
	log "github.com/erdedan1/shared/logger"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	address string
	log     log.Logger
	server  *grpc.Server
	lis     net.Listener
}

func NewGRPCServer(address string, orderService usecase.OrderService, logger log.Logger, tp trace.TracerProvider, cfg config.InfrastructureConfig) (*GRPCServer, *errs.CustomError) {
	rateLimiter := newGRPCRateLimiter(
		cfg.RateLimiter.RequestsPerSecond,
		cfg.RateLimiter.Burst,
	)

	cycleBreaker := newGRPCCircuitBreaker(
		uint32(cfg.CircuitBreaker.ConsecutiveFailures),
		cfg.CircuitBreaker.HalfOpenRequests,
		cfg.CircuitBreaker.OpenTimeout,
	)

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestid.XRequestIDServerInterceptor(),
			pbLogger.LoggerServerInterceptor(logger),
			recovery.RecoveryServerInterceptor(logger),
			rateLimiter.Unary(),
			cycleBreaker.Unary(),
		),
		grpc.ChainStreamInterceptor(
			rateLimiter.Stream(),
			cycleBreaker.Stream(),
		),
	)

	handler := New(orderService, logger, tp)
	pbOrder.RegisterOrderServiceServer(server, handler)

	return &GRPCServer{
		address: address,
		log:     logger,
		server:  server,
	}, nil
}

func (s *GRPCServer) Start() *errs.CustomError {
	const method = "GRPCServer.Start"

	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		s.log.Error("GRPCServer", method, "failed to listen", err)
		return errs.New(errs.INTERNAL, "failed to listen", err)
	}
	s.lis = lis

	err = s.server.Serve(lis)
	if err != nil {
		s.log.Error("GRPCServer", method, "grpc serve error", err)
		return errs.New(errs.INTERNAL, "grpc serve error", err)
	}

	return nil
}

func (s *GRPCServer) Stop() {
	const method = "GRPCServer.Stop"
	if s.server == nil {
		return
	}
	s.server.GracefulStop()
	s.log.Info("GRPCServer", method, "grpc server stopped gracefully")
}

func IsExpectedStop(err *errs.CustomError) bool {
	return errors.Is(err, grpc.ErrServerStopped)
}
