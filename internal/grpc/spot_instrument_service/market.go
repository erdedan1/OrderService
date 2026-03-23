package spot_instrument_service

import (
	"context"

	"OrderService/config"
	"OrderService/internal/grpc/spot_instrument_service/mapper"
	"OrderService/internal/model"
	"OrderService/internal/usecase"
	grpc_client "OrderService/pkg/client/grpc"

	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen"
	"github.com/erdedan1/shared/errs"
	"go.opentelemetry.io/otel/trace"
)

type marketService struct {
	client pb.MarketServiceClient
	conn   grpc_client.IGRPCClient
	trace  trace.Tracer
}

func NewMarketService(cfg *config.Config, tp trace.TracerProvider) (*marketService, *errs.CustomError) {
	conn, err := SetupSpotInstrumentClient(cfg)
	if err != nil {
		return nil, errs.New(errs.UNAVAILABLE, err.Error(), err)
	}

	return &marketService{
		client: pb.NewMarketServiceClient(conn),
		conn:   conn,
		trace:  tp.Tracer("order-service/MarketService"),
	}, nil
}

func (s *marketService) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

func (s *marketService) ViewMarketsByRoles(ctx context.Context, request *usecase.ViewMarketsByRolesInput) ([]model.Market, *errs.CustomError) {
	resp, err := s.client.ViewMarketsByRoles(ctx, mapper.ViewMarketsRequestToProto(request))
	if err != nil {
		return nil, errs.New(errs.UNAVAILABLE, err.Error(), err)
	}

	marketsResp := make([]model.Market, 0, len(resp.Markets))
	for _, m := range resp.Markets {
		dtoM, err := mapper.ViewMarketsResponseFromProto(m)
		if err != nil {
			return nil, err
		}
		marketsResp = append(marketsResp, *dtoM)
	}

	return marketsResp, nil
}
