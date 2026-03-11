package usecase

import (
	"time"

	errs "OrderService/internal/errors"
	"OrderService/internal/model"

	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateOrderInput struct {
	MarketID  uuid.UUID
	UserID    uuid.UUID
	OrderType string
	Price     string
	UserRoles []string
	Quantity  int64
}

type CreateOrderOutput struct {
	ID     uuid.UUID
	Status string
}

type GetOrderStatusInput struct {
	UserID  uuid.UUID
	OrderID uuid.UUID
}

type GetOrderStatusOutput struct {
	Status    string
	UpdatedAt *time.Time
}

type ViewMarketsByRolesInput struct {
	UserRoles []string
}

func CreateOrderRequestToModel(request *CreateOrderInput) (*model.Order, *errors.CustomError) {
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
