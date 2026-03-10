package usecase

import (
	"context"

	"OrderService/internal/model"

	errors "github.com/erdedan1/shared/errs"
)

type MarketService interface {
	ViewMarketsByRoles(ctx context.Context, role []string) ([]model.Market, *errors.CustomError)
}
