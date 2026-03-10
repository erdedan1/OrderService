package order_service

import (
	"context"

	"OrderService/internal/grpc/order_service/mapper"
	"OrderService/internal/usecase"

	pb "github.com/erdedan1/protocol/proto/order_service/gen"
	log "github.com/erdedan1/shared/logger"
	"go.opentelemetry.io/otel"
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

func New(orderService usecase.OrderService, log log.Logger) *Handler {
	return &Handler{
		orderService: orderService,
		log:          log,
		tracer:       otel.Tracer("order-service/MarketHandler"),
	}
}

const layer = "MarketHandler"

func (h *Handler) CreateOrder(ctx context.Context, request *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	const method = "CreateOrder"

	ctx, span := h.tracer.Start(ctx, "MarketHandler.CreateOrder")
	defer span.End()

	requestDto, err := mapper.CreateOrderRequestFromProto(request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		h.log.Error(
			layer, method,
			err.Error(), err,
		)
		return nil, status.Error(grpc_codes.Code(err.Code), err.Message)
	}

	order, err := h.orderService.CreateOrder(ctx, requestDto)
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

	return mapper.CreateOrderResponseToProto(order), nil
}

func (h *Handler) GetOrderStatus(ctx context.Context, request *pb.GetOrderStatusRequest) (*pb.GetOrderStatusResponse, error) {
	const method = "GetOrderStatus"

	ctx, span := h.tracer.Start(ctx, "MarketHandler.GetOrderStatus")
	defer span.End()

	requestDto, err := mapper.GetOrderStatusRequestFromProto(request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Message)

		h.log.Error(
			layer, method,
			err.Error(), err,
		)
		return nil, status.Error(grpc_codes.Code(err.Code), err.Message)
	}

	order, err := h.orderService.GetOrderStatus(ctx, requestDto)
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

	return mapper.GetOrderStatusResponseToProto(order), nil
}

func (h *Handler) SubscribeOrderStatus(request *pb.GetOrderStatusRequest, stream pb.OrderService_SubscribeOrderStatusServer) error {
	const method = "SubscribeOrderStatus"

	ctx := stream.Context()

	ctx, span := h.tracer.Start(ctx, "MarketHandler.SubscribeOrderStatus")
	defer span.End()

	requestDto, err := mapper.GetOrderStatusRequestFromProto(request)
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

	ch, err := h.orderService.SubscribeOrderStatus(ctx, requestDto)
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

			err := stream.Send(mapper.GetOrderStatusResponseToProto(order))
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
