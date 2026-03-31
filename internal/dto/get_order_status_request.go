package dto

import (
	errs "OrderService/internal/errors"

	pb "github.com/erdedan1/protocol/proto/order_service/gen/v1"
	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
)

type GetOrderStatusRequest struct {
	UserUUID  uuid.UUID
	OrderUUID uuid.UUID
}

func (g *GetOrderStatusRequest) FromProto(request *pb.GetOrderStatusRequest) (*GetOrderStatusRequest, *errors.CustomError) {
	if err := request.Validate(); err != nil {
		return nil, errs.ErrInvalidArgument
	}

	userId, _ := uuid.Parse(request.UserUuid)
	orderId, _ := uuid.Parse(request.OrderUuid)

	g.UserUUID = userId
	g.OrderUUID = orderId

	return g, nil
}
