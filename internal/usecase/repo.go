package usecase

import (
	"context"
	"time"

	"OrderService/internal/dto"
	"OrderService/internal/model"

	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
)

//go:generate mockery --name=OrderRepo --output=../../mocks --outpkg=mocks
type OrderRepo interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, *errors.CustomError)
	GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, *errors.CustomError)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, order model.OrderStatus) *errors.CustomError
}

//go:generate mockery --name=UserRepo --output=../../mocks --outpkg=mocks
type UserRepo interface {
	CreateUser(ctx context.Context, user model.User)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, *errors.CustomError)
}

//go:generate mockery --name=MarketCacheRepo --output=../../mocks --outpkg=mocks
type MarketCacheRepo interface {
	Set(ctx context.Context, key string, value []dto.ViewMarketsResponse, ttl time.Duration) *errors.CustomError
	Get(ctx context.Context, key string) ([]dto.ViewMarketsResponse, *errors.CustomError)
	Del(ctx context.Context, key string) *errors.CustomError
}

//go:generate mockery --name=OrderStatusSubscriber --output=../../mocks --outpkg=mocks
type OrderStatusSubscriber interface {
	SubscribeOrderStatus(ctx context.Context, orderID uuid.UUID) (<-chan model.OrderStatus, *errors.CustomError)
}

//go:generate mockery --name=OrderStatusPublisher --output=../../mocks --outpkg=mocks
type OrderStatusPublisher interface {
	PublishOrderStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) *errors.CustomError
}
