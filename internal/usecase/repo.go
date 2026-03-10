package usecase

import (
	"context"
	"time"

	"OrderService/internal/dto"
	"OrderService/internal/model"

	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, order *model.Order) (*model.Order, *errors.CustomError)
	GetOrder(ctx context.Context, id uuid.UUID) (*model.Order, *errors.CustomError)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, order model.OrderStatus) *errors.CustomError
}

type UserRepo interface {
	CreateUser(ctx context.Context, user model.User)
	GetUserById(ctx context.Context, id uuid.UUID) (*model.User, *errors.CustomError)
}

type MarketCacheRepo interface {
	Set(ctx context.Context, key string, value []dto.ViewMarketsResponse, ttl time.Duration) *errors.CustomError
	Get(ctx context.Context, key string) ([]dto.ViewMarketsResponse, *errors.CustomError)
	Del(ctx context.Context, key string) *errors.CustomError
}
