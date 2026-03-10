package order_status

import (
	"context"

	errs "OrderService/internal/errors"
	"OrderService/internal/model"
	"OrderService/pkg/cache"

	errors "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RedisPublisher struct {
	client cache.RedisClient
	log    log.Logger
	tracer trace.Tracer
}

func NewRedisPublisher(client cache.RedisClient, logger log.Logger) *RedisPublisher {
	return &RedisPublisher{
		client: client,
		log:    logger,
		tracer: otel.Tracer("order-service/RedisOrderStatusPublisher"),
	}
}

const publisherLayer = "RedisOrderStatusPublisher"

func (p *RedisPublisher) PublishOrderStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) *errors.CustomError {
	const method = "PublishOrderStatus"

	ctx, span := p.tracer.Start(ctx, "OrderStatusPublisher.PublishOrderStatus")
	defer span.End()

	channelName := orderStatusChannel(orderID)
	payload := string(status)

	span.SetAttributes(
		attribute.String("order.id", orderID.String()),
		attribute.String("channel", channelName),
		attribute.String("status", payload),
	)

	if err := p.client.Publish(ctx, channelName, payload).Err(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, errs.ErrUnavailableRedis.Message)

		p.log.Error(publisherLayer, method, err.Error(), err, "order_id", orderID, "status", payload)
		return errs.ErrUnavailableRedis
	}

	span.SetStatus(codes.Ok, "order status published")
	return nil
}
