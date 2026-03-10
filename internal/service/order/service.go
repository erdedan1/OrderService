package order

import (
	"context"
	"time"

	"OrderService/internal/dto"
	errs "OrderService/internal/errors"
	"OrderService/internal/model"
	"OrderService/internal/usecase"

	errors "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	repo   *usecase.Repositories
	pb     usecase.GRPCServices
	l      log.Logger
	tracer trace.Tracer
}

func New(repo *usecase.Repositories, log log.Logger, pb usecase.GRPCServices) *Service {
	return &Service{
		repo:   repo,
		pb:     pb,
		l:      log,
		tracer: otel.Tracer("order-service/Service"),
	}
}

const layer = "OrderService"

func (s *Service) CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, *errors.CustomError) {
	const method = "CreateOrder"

	ctx, span := s.tracer.Start(ctx, "OrderService.CreateOrder")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", request.UserID.String()),
	)

	user, err := s.repo.UserRepo.GetUserById(ctx, request.UserID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.l.Error(
			layer, method,
			err.Error(), err,
			"user_id", request.UserID,
		)
		return nil, err
	}
	if !user.CheckRoles(request.UserRoles) {
		span.RecordError(err)
		span.SetStatus(codes.Error, "no access rights")

		s.l.Error(
			layer, method,
			"user has no acces to market",
			errs.ErrUserHasNoAccessToMarket,
			"user_id", request.UserID,
		)
		return nil, errs.ErrUserHasNoAccessToMarket
	}

	cacheKey := "markets:" + request.UserID.String()
	marketsCache, err := s.repo.MarketCache.Get(ctx, cacheKey)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.l.Error(
			layer, method,
			err.Error(), err,
		)
		return nil, err
	}
	if len(marketsCache) == 0 || marketsCache == nil {
		markets, err := s.pb.MarketService.ViewMarketsByRoles(ctx, dto.NewViewMarketsRequestFromRoles(user.Roles))
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			s.l.Error(layer, method,
				err.Error(), err,
				"user_id", request.UserID,
			)
			return nil, err
		}
		if len(markets) == 0 {
			span.RecordError(err)
			span.SetStatus(codes.Error, "not found markets")

			s.l.Error(
				layer, method,
				errs.ErrMarketNotFound.Message,
				errs.ErrMarketNotFound,
				"user_id", request.UserID,
			)
			return nil, errs.ErrMarketNotFound
		}
		err = s.repo.MarketCache.Set(ctx, cacheKey, markets, 5*time.Minute)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			s.l.Error(
				layer, method,
				err.Message, err,
				"user_id", request.UserID,
			)
			return nil, err
		}
	}

	req, err := createOrderRequestToModel(request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.l.Error(
			layer, method,
			err.Error(), err,
			"user_id", request.UserID,
		)
		return nil, err
	}

	order, err := s.repo.OrderRepo.CreateOrder(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.l.Error(
			layer, method,
			err.Error(), err,
			"user_id", request.UserID,
		)
		return nil, err
	}

	span.SetStatus(codes.Ok, "order success created")

	s.l.Debug(layer, method, "order success created")

	return &dto.CreateOrderResponse{
		ID:     order.ID,
		Status: order.Status.ToString(),
	}, nil
}

func (s *Service) GetOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (*dto.GetOrderStatusResponse, *errors.CustomError) {
	const method = "GetOrderStatus"

	ctx, span := s.tracer.Start(ctx, "OrderService.GetOrderStatus")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", request.UserID.String()),
		attribute.String("order.id", request.OrderID.String()),
	)

	order, err := s.repo.OrderRepo.GetOrder(ctx, request.OrderID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.l.Error(
			layer, method,
			err.Error(), err,
			"user_id", request.UserID,
			"order_id", request.OrderID,
		)
		return nil, err
	}

	if order.UserID != request.UserID {
		span.RecordError(errs.ErrInvalidUserID)
		span.SetStatus(codes.Error, errs.ErrInvalidUserID.Message)

		s.l.Error(
			layer,
			method,
			errs.ErrInvalidUserID.Message,
			errs.ErrInvalidUserID,
			"user_id", request.UserID,
			"order_id", request.OrderID,
		)
		return nil, errs.ErrInvalidUserID
	}

	span.SetStatus(codes.Ok, "get order success")

	s.l.Debug(
		layer,
		method,
		"get order info",
		"order_id", order.ID,
		"status", order.Status,
	)

	return &dto.GetOrderStatusResponse{Status: order.Status.ToString()}, nil
}

func (s *Service) SubscribeOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (<-chan *dto.GetOrderStatusResponse, *errors.CustomError) {
	const method = "SubscribeOrderStatus"

	ctx, span := s.tracer.Start(ctx, "OrderService.SubscribeOrderStatus")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", request.UserID.String()),
		attribute.String("order.id", request.OrderID.String()),
	)

	ch := make(chan *dto.GetOrderStatusResponse)

	order, err := s.repo.OrderRepo.GetOrder(ctx, request.OrderID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.l.Error(
			layer, method,
			err.Error(), err,
			"order_id", request.OrderID,
		)
		return nil, err
	}

	if order.UserID != request.UserID {
		defer close(ch)
		span.RecordError(errs.ErrInvalidUserID)
		span.SetStatus(codes.Error, errs.ErrInvalidUserID.Message)

		s.l.Error(
			layer, method,
			errs.ErrInvalidUserID.Message,
			errs.ErrInvalidUserID,
			"order_id", request.OrderID,
			"user_id", request.UserID,
		)
		return ch, errs.ErrInvalidUserID
	}

	if order.Status == model.StatusClosed {
		defer close(ch)

		span.SetStatus(codes.Ok, "order already completed and closed")

		s.l.Debug(
			layer, method,
			"order already completed and closed",
			"user_id", request.UserID,
			"order_id", request.OrderID,
		)
		ch <- &dto.GetOrderStatusResponse{
			Status:    order.Status.ToString(),
			UpdatedAt: &order.UpdateAt,
		}
		return ch, nil
	}

	go func() {
		defer close(ch)

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		idx := 0

		for {
			select {
			case <-ctx.Done():
				return
				//todo убрать все сделать другое
			case <-ticker.C:
				if idx >= len(model.OrderStatusProcessing) {
					return
				}
				order.Status = model.OrderStatusCreated
				order.UpdateAt = time.Now()
				idx++
				err = s.repo.OrderRepo.UpdateOrderStatus(ctx, order.ID, order.Status)
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Message)

					s.l.Error(
						layer, method,
						err.Error(), err,
						"user_id", request.UserID,
						"order_id", request.OrderID,
					)
					return
				}
				s.l.Debug(
					layer, method,
					"update new order status",
					"user_id", request.UserID,
					"order_id", request.OrderID,
				)

				ch <- &dto.GetOrderStatusResponse{
					Status:    order.Status.ToString(),
					UpdatedAt: &order.UpdateAt,
				}
			}
		}
	}()
	span.SetStatus(codes.Ok, "get order success")
	return ch, nil
}

func createOrderRequestToModel(request *dto.CreateOrderRequest) (*model.Order, *errors.CustomError) {
	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		return nil, errs.ErrInvalidArgument
	}

	return &model.Order{
		UserID:    request.UserID,
		MarketID:  request.MarketID,
		Quantity:  request.Quantity,
		Type:      request.OrderType,
		Status:    model.StatusCreated,
		Price:     price,
		CreatedAt: time.Now(),
	}, nil
}
