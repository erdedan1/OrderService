package spot_instrument_service

import (
	"OrderService/config"
	"OrderService/internal/dto"
	grpc_client "OrderService/pkg/client/grpc"
	"context"

	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen"
	errors "github.com/erdedan1/shared/errs"
)

type IMarketDriver interface {
	ViewMarketsByRoles(ctx context.Context, req *dto.ViewMarketsRequest) ([]dto.ViewMarketsResponse, *errors.CustomError)
}

type Driver interface {
	IMarketDriver
}

type driverImpl struct { //поменять
	IMarketDriver
}

func New(
	cfg *config.Config,
) (*driverImpl, []grpc_client.IGRPCClient, error) {
	grpcConns := make([]grpc_client.IGRPCClient, 0)
	spotInstrumentConn, err := SetupSpotInstrumentClient(cfg)
	if err != nil {
		return nil, nil, err
	}
	grpcConns = append(grpcConns, spotInstrumentConn)

	return &driverImpl{
		IMarketDriver: NewMarketDriver(
			pb.NewMarketServiceClient(spotInstrumentConn),
		),
	}, grpcConns, nil
}
