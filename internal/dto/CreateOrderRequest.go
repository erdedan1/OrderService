package dto

import (
	"time"

	errs "OrderService/internal/errors"
	"OrderService/internal/model"

	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateOrderRequest struct {
	UserID    uuid.UUID
	MarketID  uuid.UUID
	OrderType string
	Price     string
	UserRoles []string
	Quantity  int64
}

func CreateOrderRequestToModel(request *CreateOrderRequest) (*model.Order, *errors.CustomError) {
	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	return &model.Order{
		UserID:    request.UserID,
		MarketID:  request.MarketID,
		Quantity:  request.Quantity,
		Type:      request.OrderType,
		Status:    model.StatusCreated,
		Price:     price,
		CreatedAt: time.Now(),
	}, nil
}
