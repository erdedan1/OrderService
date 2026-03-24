package dto

import (
	pb "github.com/erdedan1/protocol/proto/order_service/gen/v2"
	errors "github.com/erdedan1/shared/errs"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	UserID    uuid.UUID
	MarketID  uuid.UUID
	OrderType string
	UserRoles []string
	Price     int64
	Quantity  int64
}

func (c *CreateOrderRequest) FromProto(request *pb.CreateOrderRequest) (*CreateOrderRequest, *errors.CustomError) {
	if err := request.Validate(); err != nil {
		return nil, errors.New(errors.INVALID_ARGUMENT, "invalid request", err)
	}

	userId, _ := uuid.Parse(request.UserId)
	marketId, _ := uuid.Parse(request.MarketId)

	c.UserID = userId
	c.MarketID = marketId
	c.OrderType = request.OrderType
	c.Price = request.Price
	c.UserRoles = m.UserRolesFromProto(request.UserRoles)
	c.Quantity = request.Quantity

	return c, nil
}
