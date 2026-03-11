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
	errors "github.com/erdedan1/shared/errs"
)

type marketService struct {
	client pb.MarketServiceClient
	conn   grpc_client.IGRPCClient
}

func NewMarketService(cfg *config.Config) (*marketService, error) {
	conn, err := SetupSpotInstrumentClient(cfg)
	if err != nil {
		return nil, err
	}

	return &marketService{
		client: pb.NewMarketServiceClient(conn),
		conn:   conn,
	}, nil
}

func (s *marketService) Close() error {
	if s.conn == nil {
		return nil
	}

	return s.conn.Close()
}

func (s *marketService) ViewMarketsByRoles(ctx context.Context, request *usecase.ViewMarketsByRolesInput) ([]model.Market, *errors.CustomError) {
	resp, err := s.client.ViewMarketsByRoles(ctx, mapper.ViewMarketsRequestToProto(request))
	if err != nil {
		return nil, errs.New(errs.UNAVAILABLE, err.Error())
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
