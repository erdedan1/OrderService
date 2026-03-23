package order

import (
	"OrderService/internal/dto"
	errs "OrderService/internal/errors"
	"context"

	errors "github.com/erdedan1/shared/errs"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

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

	return &dto.GetOrderStatusResponse{
		Status:    string(order.Status),
		UpdatedAt: &order.UpdatedAt,
	}, nil
}
