package order_service

import (
	"context"

	"OrderService/internal/dto"
	"OrderService/internal/usecase"

	pb "github.com/erdedan1/protocol/proto/order_service/gen/v1"
	log "github.com/erdedan1/shared/logger"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	orderService usecase.OrderService
	log          log.Logger
	tracer       trace.Tracer

	pb.UnimplementedOrderServiceServer
}

func New(orderService usecase.OrderService, log log.Logger, tp trace.TracerProvider) *Handler {
	return &Handler{
		orderService: orderService,
		log:          log,
		tracer:       tp.Tracer("order-service/OrderHandler"),
	}
}

const layer = "OrderHandler"

func (h *Handler) CreateOrder(ctx context.Context, request *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	const method = "CreateOrder"

	ctx, span := h.tracer.Start(ctx, "OrderHandler.CreateOrder")
	defer span.End()

	if !checkUser(ctx, request.UserUuid) {
		return nil, status.Error(3, "invalid user-uuid")
	}

	dto, err := new(dto.CreateOrderRequest).FromProto(request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		h.log.Error(
			layer, method,
			err.Error(), err,
		)
		return nil, status.Error(grpc_codes.Code(err.Code), err.Message)
	}

	order, err := h.orderService.CreateOrder(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		h.log.Error(
			layer, method,
			err.Error(), err,
		)
		return nil, status.Error(grpc_codes.Code(err.Code), err.Message)
	}

	span.SetStatus(codes.Ok, "order success created")

	return order.ToProto(), nil
}

func (h *Handler) GetOrderStatus(ctx context.Context, request *pb.GetOrderStatusRequest) (*pb.GetOrderStatusResponse, error) {
	const method = "GetOrderStatus"

	ctx, span := h.tracer.Start(ctx, "OrderHandler.GetOrderStatus")
	defer span.End()

	if !checkUser(ctx, request.UserUuid) {
		return nil, status.Error(3, "invalid user-uuid")
	}

	dto, err := new(dto.GetOrderStatusRequest).FromProto(request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		h.log.Error(
			layer, method,
			err.Error(), err,
		)
		return nil, status.Error(grpc_codes.Code(err.Code), err.Message)
	}

	order, err := h.orderService.GetOrderStatus(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		h.log.Error(
			layer, method,
			err.Error(), err,
		)
		return nil, status.Error(grpc_codes.Code(err.Code), err.Message)
	}

	span.SetStatus(codes.Ok, "get order success")

	return order.ToProto(), nil
}

func (h *Handler) SubscribeOrderStatus(request *pb.GetOrderStatusRequest, stream pb.OrderService_SubscribeOrderStatusServer) error {
	const method = "SubscribeOrderStatus"

	ctx := stream.Context()

	if !checkUser(ctx, request.UserUuid) {
		return status.Error(3, "invalid user-uuid")
	}

	ctx, span := h.tracer.Start(ctx, "OrderHandler.SubscribeOrderStatus")
	defer span.End()

	dto, err := new(dto.GetOrderStatusRequest).FromProto(request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		h.log.Error(
			layer, method,
			err.Error(), err,
		)
		return status.Error(grpc_codes.Code(err.Code), err.Message)
	}

	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	ch, err := h.orderService.SubscribeOrderStatus(ctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		h.log.Error(
			layer, method,
			err.Error(), err,
		)
		return status.Error(grpc_codes.Code(err.Code), err.Message)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case order, ok := <-ch:
			if !ok {
				return nil
			}

			err := stream.Send(order.ToProto())
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				h.log.Error(
					layer, method,
					err.Error(), err,
				)
				return err
			}
		}
	}
}

func checkUser(ctx context.Context, requestUserUUID string) bool {
	userUUID := clientKeyFromContext(ctx)

	if userUUID != requestUserUUID {
		return false
	}
	return true
}
