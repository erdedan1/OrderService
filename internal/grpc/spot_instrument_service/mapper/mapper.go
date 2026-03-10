package mapper

import (
	"OrderService/internal/dto"
	errs "OrderService/internal/errors"

	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen"
	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ViewMarketsRequestToProto(req *dto.ViewMarketsRequest) *pb.ViewMarketsRequest {
	result := &pb.ViewMarketsRequest{UserRoles: make([]string, 0, len(req.UserRoles))}
	result.UserRoles = append(result.UserRoles, req.UserRoles...)
	return result
}

func ViewMarketsResponseFromProto(market *pb.Market) (*dto.ViewMarketsResponse, *errors.CustomError) {
	id, err := uuid.Parse(market.Id)
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	return &dto.ViewMarketsResponse{
		ID:        id,
		Name:      market.Name,
		Enabled:   market.Enabled,
		CreatedAt: new(market.CreatedAt.AsTime()),
		UpdatedAt: new(market.UpdatedAt.AsTime()),
		DeletedAt: new(market.DeletedAt.AsTime()),
	}, nil
}

func ViewMarketsResponseToProto(resp *dto.ViewMarketsResponse) *pb.Market {
	return &pb.Market{
		Id:        resp.ID.String(),
		Name:      resp.Name,
		Enabled:   resp.Enabled,
		CreatedAt: timestamppb.New(*resp.CreatedAt),
		UpdatedAt: timestamppb.New(*resp.UpdatedAt),
		DeletedAt: timestamppb.New(*resp.DeletedAt),
	}
}
