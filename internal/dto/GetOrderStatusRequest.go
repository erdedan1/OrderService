package dto

import (
	errs "OrderService/internal/errors"

	pb "github.com/erdedan1/protocol/proto/order_service/gen"
	errors "github.com/erdedan1/shared/errs"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
)

type GetOrderStatusRequest struct {
	UserId  uuid.UUID
	OrderId uuid.UUID
}

func NewGetOrderStatusRequest(request *pb.GetOrderStatusRequest) (*GetOrderStatusRequest, *errors.CustomError) {
	err := uuid.Validate(request.UserId)
	if err != nil {
		return nil, errs.ErrIvalidArgument
	}

	err = uuid.Validate(request.OrderId)
	if err != nil {
		return nil, errs.ErrIvalidArgument
	}

	return &GetOrderStatusRequest{
		UserId:  *m.FromUUIDProto(request.UserId),
		OrderId: *m.FromUUIDProto(request.OrderId),
	}, nil
}
