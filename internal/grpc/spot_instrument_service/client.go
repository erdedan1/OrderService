package spot_instrument_service

import (
	"OrderService/config"
	grpc_client "OrderService/pkg/client/grpc"

	requestid "github.com/erdedan1/shared_for_homework/pkg/interceptors/request_id"
	"google.golang.org/grpc"
)

func SetupSpotInstrumentClient(
	cfg *config.Config,
) (grpc_client.IGRPCClient, error) {
	conn, err := grpc_client.New(
		cfg.GRPCApi.SpotInstrumentServiceHost,
		cfg,
		grpc.WithUnaryInterceptor(requestid.XRequestIDClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
