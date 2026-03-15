package usecase

import (
	"context"

	errors "github.com/erdedan1/shared/errs"
)

//go:generate mockery --name=OrderService --output=../../mocks --outpkg=mocks
type OrderService interface {
	CreateOrder(ctx context.Context, request *CreateOrderInput) (*CreateOrderOutput, *errors.CustomError)
	GetOrderStatus(ctx context.Context, request *GetOrderStatusInput) (*GetOrderStatusOutput, *errors.CustomError)
	SubscribeOrderStatus(ctx context.Context, request *GetOrderStatusInput) (<-chan *GetOrderStatusOutput, *errors.CustomError)
}
