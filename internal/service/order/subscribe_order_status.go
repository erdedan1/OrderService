package order

import (
	"context"
	"fmt"
	"time"

	"OrderService/config"
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
		attribute.String("user.id", request.UserID.String()),
		attribute.String("order.id", request.OrderID.String()),
	)

	ch := make(chan *dto.GetOrderStatusResponse, 1)

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

		ch <- &dto.GetOrderStatusResponse{Status: order.Status.ToString(), UpdatedAt: order.UpdatedAt}

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

	go func(initialStatus model.OrderStatus, initialUpdatedAt *time.Time) {
		defer close(ch)

		ch <- &dto.GetOrderStatusResponse{Status: initialStatus.ToString(), UpdatedAt: initialUpdatedAt}
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

func (s *Service) publishOrderLifecycle(ctx context.Context, orderID uuid.UUID, initialStatus model.OrderStatus) {
	const method = "publishOrderLifecycle"

	s.log.Debug(layer, method, "start new sobitie")

	status := initialStatus
	for {
		nextStatus, hasNext := model.NextOrderStatus(status)
		fmt.Println(nextStatus, "1")
		if !hasNext {
			fmt.Println(nextStatus, "2")
			return
		}

		select {
		case <-ctx.Done():
			fmt.Println(nextStatus, "3")
			return
		case <-time.After(config.Global.Infrastructure.OrderLifecycleConfig.StepInterval):
			fmt.Println(nextStatus, "4")
		}
		fmt.Println(nextStatus, "5")
		if updateErr := s.UpdateOrderStatus(ctx, orderID, nextStatus); updateErr != nil {
			fmt.Println(nextStatus, "6")
			return
		}
		fmt.Println(nextStatus, "7")
		status = nextStatus
	}
}
