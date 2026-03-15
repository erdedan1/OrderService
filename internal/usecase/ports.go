package usecase

import (
	"context"

	"OrderService/internal/model"

	errors "github.com/erdedan1/shared/errs"
)

//go:generate mockery --name=MarketService --output=../../mocks --outpkg=mocks
type MarketService interface {
	ViewMarketsByRoles(ctx context.Context, request *ViewMarketsByRolesInput) ([]model.Market, *errors.CustomError)
}
