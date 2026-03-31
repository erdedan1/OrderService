package order

import (
	"OrderService/internal/dto"
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
		attribute.String("user.id", request.UserUUID.String()),
		attribute.String("order.id", request.OrderUUID.String()),
	)

	order, err := s.orderRepo.GetOrder(ctx, request.OrderUUID, request.UserUUID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Error(
			layer, method,
			err.Error(), err,
			"user_id", request.UserUUID,
			"order_id", request.OrderUUID,
		)
		return nil, err
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
		UpdatedAt: order.UpdatedAt,
	}, nil
}
