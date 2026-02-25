package usecase

import (
	"context"

	"OrderService/internal/dto"

	errors "github.com/erdedan1/shared/errs"
)

type MarketService interface {
	ViewMarketsByRoles(ctx context.Context, req *dto.ViewMarketsRequest) ([]dto.ViewMarketsResponse, *errors.CustomError)
}

type GRPCServices struct {
	MarketService MarketService
}

func NewGRPCServices(marketService MarketService) *GRPCServices {
	return &GRPCServices{
		MarketService: marketService,
	}
}
