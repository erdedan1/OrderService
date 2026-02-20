package dto

import (
	"time"

	errs "OrderService/internal/errors"

	pb "github.com/erdedan1/protocol/proto/spot_instrument_service/gen"
	errors "github.com/erdedan1/shared/errs"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ViewMarketsResponse struct {
	ID        uuid.UUID
	Name      string
	Enabled   bool
	CreatedAt *time.Time
	UpdateAt  *time.Time
	DeletedAt *time.Time
}

func (vmr *ViewMarketsResponse) DtoToProto() *pb.Market {
	return &pb.Market{
		Id:        vmr.ID.String(),
		Name:      vmr.Name,
		Enabled:   vmr.Enabled,
		CreatedAt: timestamppb.New(*vmr.CreatedAt),
		UpdatedAt: timestamppb.New(*vmr.UpdateAt),
		DeletedAt: timestamppb.New(*vmr.DeletedAt),
	}
}

func NewViewMarketsResponse(market *pb.Market) (*ViewMarketsResponse, *errors.CustomError) {
	err := uuid.Validate(market.Id)
	if err != nil {
		return nil, errs.ErrIvalidArgument
	}

	return &ViewMarketsResponse{
		ID:        *m.FromUUIDProto(market.Id),
		Name:      market.Name,
		Enabled:   market.Enabled,
		CreatedAt: m.ToDtoTime(market.CreatedAt),
		UpdateAt:  m.ToDtoTime(market.UpdatedAt),
		DeletedAt: m.ToDtoTime(market.DeletedAt),
	}, nil
}
