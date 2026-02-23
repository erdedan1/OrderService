package order

import (
	"context"
	"fmt"
	"time"

	"OrderService/internal/dto"
	errs "OrderService/internal/errors"
	"OrderService/internal/grpc/spot_instrument_service"
	"OrderService/internal/model"
	"OrderService/internal/usecase"

	pbOrder "github.com/erdedan1/protocol/proto/order_service/gen"
	errors "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
)

type Service struct {
	repo *usecase.Repositories
	pb   spot_instrument_service.Driver
	l    log.Logger

	pbOrder.UnimplementedOrderServiceServer
}

func New(repo *usecase.Repositories, log log.Logger, pb spot_instrument_service.Driver) *Service {
	return &Service{
		repo: repo,
		pb:   pb,
		l:    log.Layer("Order.Service"),
	}
}

func (s *Service) CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.CreateOrderResponse, *errors.CustomError) {
	const method = "CreateOrder"

	user, err := s.repo.UserRepo.GetUserById(ctx, request.UserId)
	if err != nil {
		s.l.Error(method, err.Error(), err, request.UserId)
		return nil, err
	}

	if !user.CheckRoles(request.UserRoles) {
		fmt.Println(request.UserRoles)
		fmt.Println(user.Roles)
		s.l.Error(method, "user has no acces to market ", errs.ErrUserHasNoAccessToMarket, request.UserId)
		return nil, errs.ErrUserHasNoAccessToMarket
	}

	cacheKey := "markets:" + request.UserId.String()
	marketsCache, err := s.repo.MarketCache.Get(ctx, cacheKey)
	if err != nil {
		s.l.Error(method, err.Error(), err)
		return nil, err
	}
	if len(marketsCache) == 0 {
		markets, err := s.pb.ViewMarketsByRoles(ctx, new(dto.ViewMarketsRequest).UserRolesToProto(user.Roles))
		if err != nil {
			s.l.Error(method, err.Error(), err, request.UserId)
			return nil, err
		}
		if len(markets) == 0 {
			s.l.Error(method, errs.ErrMarketNotFound.Message, errs.ErrMarketNotFound, request.UserId)
			return nil, errs.ErrMarketNotFound
		}
		err = s.repo.MarketCache.Set(ctx, cacheKey, markets, 5*time.Minute)
		if err != nil {
			s.l.Error(method, err.Message, err, request.UserId)
			return nil, err
		}
	}

	req, err := request.DtoToModel()
	if err != nil {
		s.l.Error(method, err.Error(), err, request.UserId)
		return nil, err
	}

	order, err := s.repo.OrderRepo.CreateOrder(ctx, *req)
	if err != nil {
		s.l.Error(method, err.Error(), err, request.UserId)
		return nil, err
	}

	s.l.Debug(method, "order success created")

	return &dto.CreateOrderResponse{
		ID:     order.ID,
		Status: order.Status,
	}, nil
}

func (s *Service) GetOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (*dto.GetOrderStatusResponse, *errors.CustomError) {
	const method = "GetOrderStatus"

	order, err := s.repo.OrderRepo.GetOrder(ctx, request.OrderId)
	if err != nil {
		s.l.Error(method, err.Error(), err, request.UserId, request.OrderId)
		return nil, err
	}

	if order.UserId != request.UserId {
		s.l.Error(
			method,
			errs.ErrInvalidUserID.Message,
			errs.ErrInvalidUserID,
			request.OrderId,
			request.UserId,
		)
		return nil, errs.ErrInvalidUserID
	}

	s.l.Debug(method, "get order info", order.ID, order.Status)

	return &dto.GetOrderStatusResponse{Status: order.Status}, nil
}

func (s *Service) SubscribeOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (<-chan *dto.GetOrderStatusResponse, *errors.CustomError) {
	const method = "SubscribeOrderStatus"

	ch := make(chan *dto.GetOrderStatusResponse)

	order, err := s.repo.OrderRepo.GetOrder(ctx, request.OrderId)
	if err != nil {
		s.l.Error(method, err.Error(), err, request.OrderId, request.UserId)
		return nil, err
	}

	if order.UserId != request.UserId {
		defer close(ch)
		s.l.Error(
			method,
			errs.ErrInvalidUserID.Message,
			errs.ErrInvalidUserID,
			request.OrderId,
			request.UserId,
		)
		return ch, errs.ErrInvalidUserID
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

			case <-ticker.C:
				if idx >= len(model.OrderStatusProcessing) {
					return
				}
				order.Status = model.OrderStatusProcessing[idx]
				order.UpdateAt = time.Now()
				idx++
				orderUpdated, err := s.repo.OrderRepo.UpdateOrder(ctx, order.ID, *order)
				if err != nil {
					s.l.Error(method, err.Error(), err, request.UserId, request.OrderId)
					return
				}
				s.l.Debug(method, "update new order status", request.UserId, request.OrderId)
				ch <- &dto.GetOrderStatusResponse{
					Status: orderUpdated.Status,
				}
			}
		}
	}()

	return ch, nil
}
