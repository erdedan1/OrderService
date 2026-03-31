package order

import (
	"context"
	"time"

	"OrderService/internal/dto"
	errs "OrderService/internal/errors"
	"OrderService/internal/model"

	errors "github.com/erdedan1/shared/errs"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (s *Service) SubscribeOrderStatus(ctx context.Context, request *dto.GetOrderStatusRequest) (<-chan *dto.GetOrderStatusResponse, *errors.CustomError) {
	const method = "SubscribeOrderStatus"

	ctx, span := s.tracer.Start(ctx, "OrderService.SubscribeOrderStatus")
	defer span.End()

	span.SetAttributes(
		attribute.String("user.id", request.UserUUID.String()),
		attribute.String("order.id", request.OrderUUID.String()),
	)

	ch := make(chan *dto.GetOrderStatusResponse, 1)

	order, err := s.orderRepo.GetOrder(ctx, request.OrderUUID, request.UserUUID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Error(
			layer, method,
			err.Error(), err,
			"order_id", request.OrderUUID,
		)
		return nil, err
	}

	if order.Status == model.StatusClosed {
		defer close(ch)

		ch <- &dto.GetOrderStatusResponse{Status: order.Status.ToString(), UpdatedAt: order.UpdatedAt}

		span.SetStatus(codes.Ok, "order already completed and closed")

		s.log.Debug(
			layer, method,
			"order already completed and closed",
			"user_id", request.UserUUID,
			"order_id", request.OrderUUID,
		)

		return ch, nil
	}

	if s.orderStatusSubscriber == nil {
		defer close(ch)
		return ch, errs.ErrUnavailableRedis
	}

	statusCh, err := s.orderStatusSubscriber.SubscribeOrderStatus(ctx, request.OrderUUID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		s.log.Error(layer, method, err.Error(), err, "order_id", request.OrderUUID, "user_id", request.UserUUID)
		return nil, err
	}

	lifecircuitCtx, cancel := context.WithTimeout(context.Background(), s.cfg.Infrastructure.OrderLifecircuitConfig.TimeOut)
	defer cancel() //проверить таймаут

	go s.publishOrderLifecircuit(lifecircuitCtx, request.UserUUID, order.ID, order.Status)

	go func(initialStatus model.OrderStatus, initialUpdatedAt *time.Time) {
		defer close(ch)

		select {
		case <-ctx.Done():
			return
		case ch <- &dto.GetOrderStatusResponse{Status: initialStatus.ToString(), UpdatedAt: initialUpdatedAt}:
		}
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

				select {
				case <-ctx.Done():
					return
				case ch <- &dto.GetOrderStatusResponse{Status: status.ToString(), UpdatedAt: &now}:
				}

				if status == model.StatusClosed {
					return
				}
			}
		}
	}(order.Status, order.UpdatedAt)

	span.SetStatus(codes.Ok, "subscribe order status started")

	return ch, nil
}

func (s *Service) publishOrderLifecircuit(ctx context.Context, userID, orderID uuid.UUID, initialStatus model.OrderStatus) {
	const method = "publishOrderLifecircuit"

	s.log.Debug(layer, method, "start new sobitie")

	status := initialStatus
	for {
		nextStatus, hasNext := model.NextOrderStatus(status)
		if !hasNext {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(s.cfg.Infrastructure.OrderLifecircuitConfig.StepInterval):
		}
		if updateErr := s.UpdateOrderStatus(ctx, userID, orderID, nextStatus); updateErr != nil {
			return
		}
		status = nextStatus
	}
}
