package usecase

import (
	"context"

	"OrderService/internal/dto"

	errors "github.com/erdedan1/shared/errs"
)

//go:generate mockery --name=MarketService --output=../../mocks --outpkg=mocks
type MarketService interface {
	ViewMarketsByRoles(ctx context.Context, request *dto.ViewMarketsRequest) ([]dto.ViewMarketsResponse, *errors.CustomError)
}
