package dto

import (
	pb "github.com/erdedan1/protocol/proto/order_service/gen/v2"
	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
)

type GetOrderStatusRequest struct {
	UserID  uuid.UUID
	OrderID uuid.UUID
}

func (g *GetOrderStatusRequest) FromProto(request *pb.GetOrderStatusRequest) (*GetOrderStatusRequest, *errors.CustomError) {
	if err := request.Validate(); err != nil {
		return nil, errors.New(errors.INVALID_ARGUMENT, "invalid request", err)
	}

	userId, _ := uuid.Parse(request.UserId)
	orderId, _ := uuid.Parse(request.OrderId)

	g.UserID = userId
	g.OrderID = orderId

	return g, nil
}
