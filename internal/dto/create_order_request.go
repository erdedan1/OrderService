package dto

import (
	pb "github.com/erdedan1/protocol/proto/order_service/gen/v1"
	"github.com/shopspring/decimal"

	errs "OrderService/internal/errors"

	errors "github.com/erdedan1/shared/errs"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	UserUUID   uuid.UUID
	MarketUUID uuid.UUID
	OrderType  string
	UserRole   string
	Price      decimal.Decimal
	Quantity   int64
}

func (c *CreateOrderRequest) FromProto(request *pb.CreateOrderRequest) (*CreateOrderRequest, *errors.CustomError) {
	if err := request.Validate(); err != nil {
		return nil, errs.ErrInvalidArgument
	}

	userId, _ := uuid.Parse(request.UserUuid)
	marketId, _ := uuid.Parse(request.MarketUuid)
	price, err := decimal.NewFromString(request.Price.String())
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	c.UserUUID = userId
	c.MarketUUID = marketId
	c.OrderType = m.OrderTypeFromProto(request.OrderType)
	c.Price = price
	c.UserRole = m.UserRoleFromProtoOrder(request.UserRole)
	c.Quantity = request.Quantity

	return c, nil
}
