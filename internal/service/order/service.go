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
	orderRepo             usecase.OrderRepo
	userRepo              usecase.UserRepo
	marketCache           usecase.MarketCacheRepo
	marketSrv             usecase.MarketService
	orderStatusSubscriber usecase.OrderStatusSubscriber
	orderStatusPublisher  usecase.OrderStatusPublisher
	log                   log.Logger
	tracer                trace.Tracer
}

func New(
	repo usecase.OrderRepo,
	userRepo usecase.UserRepo,
	marketCache usecase.MarketCacheRepo,
	marketSrv usecase.MarketService,
	orderStatusSubscriber usecase.OrderStatusSubscriber,
	orderStatusPublisher usecase.OrderStatusPublisher,
	log log.Logger,
) *Service {
	return &Service{
		orderRepo:             repo,
		userRepo:              userRepo,
		marketCache:           marketCache,
		marketSrv:             marketSrv,
		orderStatusSubscriber: orderStatusSubscriber,
		orderStatusPublisher:  orderStatusPublisher,
		log:                   log,
		tracer:                otel.Tracer("order-service/Service"),
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

	user, err := s.userRepo.GetUserById(ctx, request.UserID)
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
	if !user.CheckRoles(request.UserRoles) {
		span.RecordError(err)
		span.SetStatus(codes.Error, "no access rights")

		s.log.Error(
			layer, method,
			"user has no acces to market",
			errs.ErrUserHasNoAccessToMarket,
			"user_id", request.UserID,
		)
		return nil, errs.ErrUserHasNoAccessToMarket
	}

	cacheKey := "markets:" + request.UserID.String()
	marketsCache, err := s.marketCache.Get(ctx, cacheKey)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Error(
			layer, method,
			err.Error(), err,
		)
		return nil, err
	}
	if len(marketsCache) == 0 || marketsCache == nil {
		markets, err := s.marketSrv.ViewMarketsByRoles(ctx, user.Roles)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			s.log.Error(layer, method,
				err.Error(), err,
				"user_id", request.UserID,
			)
			return nil, err
		}
		if len(markets) == 0 {
			span.RecordError(err)
			span.SetStatus(codes.Error, "not found markets")

			s.log.Error(
				layer, method,
				errs.ErrMarketNotFound.Message,
				errs.ErrMarketNotFound,
				"user_id", request.UserID,
			)
			return nil, errs.ErrMarketNotFound
		}
		err = s.marketCache.Set(ctx, cacheKey, markets, 5*time.Minute)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())

			s.log.Error(
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

		s.log.Error(
			layer, method,
			err.Error(), err,
			"user_id", request.UserID,
		)
		return nil, err
	}

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

	if s.orderStatusPublisher != nil {
		if publishErr := s.orderStatusPublisher.PublishOrderStatus(ctx, order.ID, order.Status); publishErr != nil {
			s.log.Error(layer, method, publishErr.Error(), publishErr, "order_id", order.ID, "status", order.Status)
		}
	}

	span.SetStatus(codes.Ok, "order success created")

	s.log.Debug(layer, method, "order success created")

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

	order, err := s.orderRepo.GetOrder(ctx, request.OrderID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Error(
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

		s.log.Error(
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

	s.log.Debug(
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

	order, err := s.orderRepo.GetOrder(ctx, request.OrderID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Error(
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

		s.log.Error(
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

		ch <- &dto.GetOrderStatusResponse{Status: order.Status.ToString(), UpdatedAt: &order.UpdatedAt}

		span.SetStatus(codes.Ok, "order already completed and closed")

		s.log.Debug(
			layer, method,
			"order already completed and closed",
			"user_id", request.UserID,
			"order_id", request.OrderID,
		)

		return ch, nil
	}

	if s.orderStatusSubscriber == nil {
		defer close(ch)
		return ch, errs.ErrUnavailableRedis
	}

	statusCh, err := s.orderStatusSubscriber.SubscribeOrderStatus(ctx, request.OrderID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		s.log.Error(layer, method, err.Error(), err, "order_id", request.OrderID, "user_id", request.UserID)
		return nil, err
	}

	go func(initialStatus model.OrderStatus, initialUpdatedAt time.Time) {
		defer close(ch)

		ch <- &dto.GetOrderStatusResponse{Status: initialStatus.ToString(), UpdatedAt: &initialUpdatedAt}
		lastStatus := initialStatus

		for {
			select {
			case <-ctx.Done():
				return
			case status, ok := <-statusCh:
				if !ok {
					return
				}

				if status == lastStatus {
					continue
				}
				lastStatus = status

				now := time.Now()
				ch <- &dto.GetOrderStatusResponse{Status: status.ToString(), UpdatedAt: &now}

				if status == model.StatusClosed {
					return
				}
			}
		}
	}(order.Status, order.UpdatedAt)

	span.SetStatus(codes.Ok, "subscribe order status started")

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
