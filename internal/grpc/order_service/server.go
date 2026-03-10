package order_service

import (
	"errors"
	"net"

	"OrderService/internal/usecase"

	pbOrder "github.com/erdedan1/protocol/proto/order_service/gen"
	pbLogger "github.com/erdedan1/shared/interceptors/logger"
	"github.com/erdedan1/shared/interceptors/recovery"
	requestid "github.com/erdedan1/shared/interceptors/request_id"
	log "github.com/erdedan1/shared/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	address string
	log     log.Logger
	server  *grpc.Server
	lis     net.Listener
}

func NewGRPCServer(address string, orderService usecase.OrderService, logger log.Logger) (*GRPCServer, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestid.XRequestIDServerInterceptor(),
			pbLogger.LoggerServerInterceptor(zapLogger),
			recovery.RecoveryServerInterceptor(zapLogger),
		),
	)

	handler := New(orderService, logger)
	pbOrder.RegisterOrderServiceServer(server, handler)

	return &GRPCServer{
		address: address,
		log:     logger,
		server:  server,
	}, nil
}

func (s *GRPCServer) Start() error {
	const method = "GRPCServer.Start"

	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		s.log.Error("GRPCServer", method, "failed to listen", err)
		return err
	}
	s.lis = lis

	err = s.server.Serve(lis)
	if err != nil {
		s.log.Error("GRPCServer", method, "grpc serve error", err)
		return err
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

func IsExpectedStop(err error) bool {
	return errors.Is(err, grpc.ErrServerStopped)
}
