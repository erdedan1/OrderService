package dto

import (
	"time"

	errs "OrderService/internal/errors"
	"OrderService/internal/model"

	pb "github.com/erdedan1/protocol/proto/order_service/gen"
	errors "github.com/erdedan1/shared/errs"
	m "github.com/erdedan1/shared/mapper"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	UserId    uuid.UUID
	MarketId  uuid.UUID
	OrderType string
	Price     string
	UserRoles []string
	Quantity  int64
}

func NewCreateOrderRequest(request *pb.CreateOrderRequest) (*CreateOrderRequest, *errors.CustomError) {
	err := uuid.Validate(request.MarketId)
	if err != nil {
		return nil, errs.ErrIvalidArgument
	}
	err = uuid.Validate(request.UserId)
	if err != nil {
		return nil, errs.ErrIvalidArgument
	}
	if request.UserRoles == nil || len(request.UserRoles) == 0 {
		return nil, errs.ErrIvalidArgument
	}
	return &CreateOrderRequest{
		MarketId:  *m.FromUUIDProto(request.MarketId),
		UserId:    *m.FromUUIDProto(request.UserId),
		OrderType: request.OrderType,
		Price:     request.Price,
		UserRoles: m.ProtoUserRolesToString(request.UserRoles),
		Quantity:  request.Quantity,
	}, nil
}

func (cor *CreateOrderRequest) DtoToModel() (*model.Order, *errors.CustomError) {
	price, err := decimal.NewFromString(cor.Price)
	if err != nil {
		return nil, errs.ErrIvalidArgument
	}
	return &model.Order{
		UserId:    cor.UserId,
		MarketId:  cor.MarketId,
		Quantity:  cor.Quantity,
		Type:      cor.OrderType,
		Status:    pb.OrderStatus_ORDER_STATUS_CREATED.String(),
		Price:     price,
		CreatedAt: time.Now(),
	}, nil
}
