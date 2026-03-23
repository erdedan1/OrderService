package dto

import (
	"time"

	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen/v2"
	errors "github.com/erdedan1/shared/errs"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ViewMarketsResponse struct {
	ID        uuid.UUID
	Name      string
	Enabled   bool
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func (v *ViewMarketsResponse) ToProto() *pb.Market {
	return &pb.Market{
		Id:        v.ID.String(),
		Name:      v.Name,
		Enabled:   v.Enabled,
		CreatedAt: timestamppb.New(*v.CreatedAt),
		UpdatedAt: timestamppb.New(*v.UpdatedAt),
		DeletedAt: timestamppb.New(*v.DeletedAt),
	}
}

func (v *ViewMarketsResponse) FromProto(market *pb.Market) (*ViewMarketsResponse, *errors.CustomError) {
	if err := market.Validate(); err != nil {
		return nil, errors.New(errors.INVALID_ARGUMENT, "invalid request", err)
	}
	id, _ := uuid.Parse(market.Id)

	v.ID = id
	v.Name = market.Name
	v.Enabled = market.Enabled
	v.CreatedAt = new(market.CreatedAt.AsTime())
	v.UpdatedAt = new(market.UpdatedAt.AsTime())
	v.DeletedAt = new(market.DeletedAt.AsTime())

	return v, nil
}
