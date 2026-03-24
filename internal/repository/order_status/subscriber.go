package order_status

import (
	"context"

	errs "OrderService/internal/errors"
	"OrderService/internal/model"
	"OrderService/pkg/cache"

	errorz "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RedisSubscriber struct {
	client cache.RedisClient
	log    log.Logger
	tracer trace.Tracer
}

func NewRedisSubscriber(client cache.RedisClient, logger log.Logger, tp trace.TracerProvider) *RedisSubscriber {
	return &RedisSubscriber{
		client: client,
		log:    logger,
		tracer: tp.Tracer("order-service/RedisOrderStatusSubscriber"),
	}
}

const layer = "RedisOrderStatusSubscriber"

func (s *RedisSubscriber) SubscribeOrderStatus(ctx context.Context, orderID uuid.UUID) (<-chan model.OrderStatus, *errorz.CustomError) {
	const method = "SubscribeOrderStatus"

	ctx, span := s.tracer.Start(ctx, "OrderStatusSubscriber.SubscribeOrderStatus")
	defer span.End()

	channelName := orderStatusChannel(orderID)
	span.SetAttributes(
		attribute.String("order.id", orderID.String()),
		attribute.String("channel", channelName),
	)

	pubsub := s.client.Subscribe(ctx, channelName)
	if _, err := pubsub.Receive(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, errs.ErrUnavailableRedis.Message)

		s.log.Error(layer, method, err.Error(), err, "order_id", orderID)
		return nil, errs.ErrUnavailableRedis
	}

	messages := pubsub.Channel()
	out := make(chan model.OrderStatus)

	go func() {
		defer close(out)
		defer func() {
			_ = pubsub.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-messages:
				if !ok {
					return
				}

				status := model.OrderStatus(msg.Payload)
				if !status.IsValid() {
					s.log.Error(layer, method, "invalid order status payload", errs.ErrInvalidArgument, "order_id", orderID, "payload", msg.Payload)
					continue
				}

				out <- status
			}
		}
	}()

	span.SetStatus(codes.Ok, "order status subscription started")
	return out, nil
}

func orderStatusChannel(orderID uuid.UUID) string {
	return "order:status:" + orderID.String()
}
