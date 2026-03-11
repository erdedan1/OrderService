package mapper

import (
	errs "OrderService/internal/errors"
	"OrderService/internal/model"
	"OrderService/internal/usecase"

	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen"
	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ViewMarketsRequestToProto(request *usecase.ViewMarketsByRolesInput) *pb.ViewMarketsRequest {
	result := &pb.ViewMarketsRequest{UserRoles: make([]string, 0, len(request.UserRoles))}
	result.UserRoles = append(result.UserRoles, request.UserRoles...)
	return result
}

func ViewMarketsResponseFromProto(market *pb.Market) (*model.Market, *errors.CustomError) {
	id, err := uuid.Parse(market.Id)
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	return &model.Market{
		ID:        id,
		Name:      market.Name,
		Enabled:   market.Enabled,
		CreatedAt: new(market.CreatedAt.AsTime()),
		UpdatedAt: new(market.UpdatedAt.AsTime()),
		DeletedAt: new(market.DeletedAt.AsTime()),
	}, nil
}

func ViewMarketsResponseToProto(resp *model.Market) *pb.Market {
	return &pb.Market{
		Id:        resp.ID.String(),
		Name:      resp.Name,
		Enabled:   resp.Enabled,
		CreatedAt: timestamppb.New(*resp.CreatedAt),
		UpdatedAt: timestamppb.New(*resp.UpdatedAt),
		DeletedAt: timestamppb.New(*resp.DeletedAt),
	}
}
