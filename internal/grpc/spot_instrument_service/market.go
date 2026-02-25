package spot_instrument_service

import (
	"OrderService/internal/dto"
	"context"

	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen"
	"github.com/erdedan1/shared/errs"
	errors "github.com/erdedan1/shared/errs"
)

type marketService struct {
	client pb.MarketServiceClient
}

func NewMarketService(client pb.MarketServiceClient) *marketService {
	return &marketService{
		client: client,
	}
}

func (s *marketService) ViewMarketsByRoles(ctx context.Context, req *dto.ViewMarketsRequest) ([]dto.ViewMarketsResponse, *errors.CustomError) {
	resp, err := s.client.ViewMarketsByRoles(ctx, req.DtoToProto())
	if err != nil {
		return nil, errs.New(errs.UNAVAILABLE, err.Error())
	}

	marketsResp := make([]dto.ViewMarketsResponse, 0, len(resp.Markets))
	for _, m := range resp.Markets {
		dtoM, err := dto.NewViewMarketsResponse(m)
		if err != nil {
			return nil, err
		}
		marketsResp = append(marketsResp, *dtoM)
	}

	return marketsResp, nil
}
