package usecase

import (
	"context"

	"OrderService/internal/dto"

	errors "github.com/erdedan1/shared/errs"
)

type MarketService interface {
	ViewMarketsByRoles(ctx context.Context, req *dto.ViewMarketsRequest) ([]dto.ViewMarketsResponse, *errors.CustomError)
}
