package market

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"OrderService/internal/dto"
	errs "OrderService/internal/errors"
	"OrderService/pkg/cache"

	error "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/go-redis/redis/v8"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type redisMarketsCache struct {
	client cache.RedisClient
	l      log.Logger
	tracer trace.Tracer
}

func NewMarketsCache(client cache.RedisClient, l log.Logger) *redisMarketsCache {
	return &redisMarketsCache{
		client: client,
		l:      l,
		tracer: otel.Tracer("order-service/RedisCacheRepo"),
	}
}

const layer = "RedisMarketCache"

func (c *redisMarketsCache) Set(ctx context.Context, key string, value []dto.ViewMarketsResponse, ttl time.Duration) *error.CustomError {
	const method = "Set"
	ctx, span := c.tracer.Start(ctx, "OrderRedisRepo.Set")
	defer span.End()

	span.SetAttributes(
		attribute.String("key", key),
	)

	data, err := json.Marshal(value)
	if err != nil {
		span.RecordError(errs.ErrFailedSerializeRedis)
		span.SetStatus(codes.Error, errs.ErrFailedSerializeRedis.Message)

		c.l.Error(layer, method, "failed to marshal data", err)
		return errs.ErrFailedSerializeRedis
	}

	err = c.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		span.RecordError(errs.ErrUnavailableRedis)
		span.SetStatus(codes.Error, errs.ErrUnavailableRedis.Message)

		c.l.Error(layer, method, "failed to set market", err)
		return errs.ErrUnavailableRedis
	}

	span.SetStatus(codes.Ok, "market success set")

	c.l.Debug(layer, method, "success set market cache")

	return nil
}

func (c *redisMarketsCache) Get(ctx context.Context, key string) ([]dto.ViewMarketsResponse, *error.CustomError) {
	const method = "Get"

	ctx, span := c.tracer.Start(ctx, "OrderRedisRepo.Get")
	defer span.End()

	span.SetAttributes(
		attribute.String("key", key),
	)

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			span.SetStatus(codes.Ok, "markets not found in redis cache")

			c.l.Debug(
				method,
				"not found market in cache",
				"error", err.Error())
			return nil, nil
		}
		span.RecordError(errs.ErrUnavailableDataRedis)
		span.SetStatus(codes.Error, errs.ErrUnavailableDataRedis.Message)

		c.l.Error(layer, method, "failed to marshal data", err)
		return nil, errs.ErrUnavailableDataRedis
	}

	var result []dto.ViewMarketsResponse
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		span.RecordError(errs.ErrFailedDeserializeRedis)
		span.SetStatus(codes.Error, errs.ErrFailedDeserializeRedis.Message)

		c.l.Error(layer, method, "failed to unmarshal data", err)
		return nil, errs.ErrFailedDeserializeRedis
	}

	span.SetStatus(codes.Ok, "market success get in cache redis")

	c.l.Debug(layer, method, "success get market in cache redis")

	return result, nil
}

func (c *redisMarketsCache) Del(ctx context.Context, key string) *error.CustomError {
	const method = "Del"

	ctx, span := c.tracer.Start(ctx, "OrderRedisRepo.Del")
	defer span.End()

	span.SetAttributes(
		attribute.String("key", key),
	)

	err := c.client.Del(ctx, key).Err()
	if err != nil {
		span.RecordError(errs.ErrDeleteRedis)
		span.SetStatus(codes.Error, errs.ErrDeleteRedis.Message)

		c.l.Error(layer, method, "failed to delete market in cache", err)
		return errs.ErrDeleteRedis
	}

	span.SetStatus(codes.Ok, "market success deleted in cache redis")

	c.l.Debug(layer, method, "success delete market in cache redis")

	return nil
}
