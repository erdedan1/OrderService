package spot_instrument_service

import (
	"OrderService/internal/dto"
	"context"

	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen"
	"github.com/erdedan1/shared/errs"
	errors "github.com/erdedan1/shared/errs"
)

type marketDriver struct {
	client pb.MarketServiceClient
}

func NewMarketDriver(client pb.MarketServiceClient) *marketDriver {
	return &marketDriver{
		client: client,
	}
}

func (s *marketDriver) ViewMarketsByRoles(
	ctx context.Context,
	req *dto.ViewMarketsRequest,
) ([]dto.ViewMarketsResponse, *errors.CustomError) {
	resp, err := s.client.ViewMarketsByRoles(ctx, req.DtoToProto())
	if err != nil {
		return nil, errs.New(errs.INTERNAL, "hui") //или errs.UNAVAILABLE
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
