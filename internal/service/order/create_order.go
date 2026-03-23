package order

import (
	"OrderService/internal/dto"
	errs "OrderService/internal/errors"
	"OrderService/internal/model"
	"context"

	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *Service) CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, *errors.CustomError) {
	const method = "CreateOrder"

	ctx, span := s.tracer.Start(ctx, "OrderService.CreateOrder")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", request.UserID.String()),
	)

	user, err := s.getAuthorizedUser(ctx, request)
	if err != nil {
		return nil, err
	}

	if err := s.ensureMarketsAccess(ctx, request.UserID, &dto.ViewMarketsRequest{UserRoles: user.Roles}); err != nil {
		return nil, err
	}

	req := model.NewOrder(
		request.UserID,
		request.MarketID,
		request.Quantity,
		request.Price,
		request.OrderType,
	)

	order, err := s.orderRepo.CreateOrder(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Error(
			layer, method,
			err.Error(), err,
			"user_id", request.UserID,
		)
		return nil, err
	}

	if publishErr := s.orderStatusPublisher.PublishOrderStatus(ctx, order.ID, order.Status); publishErr != nil {
		s.log.Error(layer, method, publishErr.Error(), publishErr, "order_id", order.ID, "status", order.Status)
	}

	span.SetStatus(codes.Ok, "order success created")
	s.log.Debug(layer, method, "order success created")

	return &dto.CreateOrderResponse{
		ID:     order.ID,
		Status: order.Status.ToString(),
	}, nil
}

func (s *Service) getAuthorizedUser(ctx context.Context, request *dto.CreateOrderRequest) (*model.User, *errors.CustomError) {
	const method = "getAuthorizedUser"

	ctx, span := s.tracer.Start(ctx, "OrderService.getAuthorizedUser")
	defer span.End()

	user, err := s.userRepo.GetUserById(ctx, request.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Error(layer, method, err.Error(), err, "user_id", request.UserID)
		return nil, err
	}

	if !user.CheckRoles(request.UserRoles) {
		span.RecordError(errs.ErrUserHasNoAccessToMarket)
		span.SetStatus(codes.Error, "no access rights")

		s.log.Error(layer, method, "user has no acces to market", errs.ErrUserHasNoAccessToMarket, "user_id", request.UserID)
		return nil, errs.ErrUserHasNoAccessToMarket
	}

	return user, nil
}

func (s *Service) ensureMarketsAccess(ctx context.Context, userID uuid.UUID, roles *dto.ViewMarketsRequest) *errors.CustomError {
	const method = "ensureMarketsAccess"

	ctx, span := s.tracer.Start(ctx, "OrderService.ensureMarketsAccess")
	defer span.End()

	cacheKey := "markets:" + userID.String()
	marketsCache, err := s.marketCache.Get(ctx, cacheKey)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Error(layer, method, err.Error(), err)
		return err
	}

	if len(marketsCache) != 0 {
		return nil
	}

	markets, err := s.marketSrv.ViewMarketsByRoles(ctx, roles)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Error(layer, method, err.Error(), err, "user_id", userID)
		return err
	}

	if len(markets) == 0 {
		span.RecordError(errs.ErrMarketNotFound)
		span.SetStatus(codes.Error, "not found markets")

		s.log.Error(layer, method, errs.ErrMarketNotFound.Message, errs.ErrMarketNotFound, "user_id", userID)
		return errs.ErrMarketNotFound
	}

	err = s.marketCache.Set(ctx, cacheKey, markets, s.cfg.Infrastructure.RedisConfig.TTL)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Error(layer, method, err.Message, err, "user_id", userID)
		return err
	}

	return nil
}
