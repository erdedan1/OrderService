package usecase

import (
	"context"

	"OrderService/internal/dto"

	errors "github.com/erdedan1/shared/errs"
)

type OrderService interface {
	CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, *errors.CustomError)
	GetOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (*dto.GetOrderStatusResponse, *errors.CustomError)
	SubscribeOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (<-chan *dto.GetOrderStatusResponse, *errors.CustomError)
}

type Services struct {
	OrderService OrderService
}

func NewServices(orderService OrderService) *Services {
	return &Services{
		OrderService: orderService,
	}
}
