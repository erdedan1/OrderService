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
)

type redisMarketsCache struct {
	client cache.RedisClient
	l      log.Logger
}

func NewMarketsCache(client cache.RedisClient, l log.Logger) *redisMarketsCache {
	return &redisMarketsCache{
		client: client,
		l:      l,
	}
}

const layer = "RedisMarketCache"

func (c *redisMarketsCache) Set(ctx context.Context, key string, value []dto.ViewMarketsResponse, ttl time.Duration) *error.CustomError {
	const method = "Set"
	data, err := json.Marshal(value)
	if err != nil {
		c.l.Error(layer, method, "failed to marshal data", err)
		return errs.ErrFailedSerializeRedis
	}

	err = c.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		c.l.Error(layer, method, "failed to set market", err)
		return errs.ErrUnavailableRedis
	}
	c.l.Debug(layer, method, "success set market cache")
	return nil
}

func (c *redisMarketsCache) Get(ctx context.Context, key string) ([]dto.ViewMarketsResponse, *error.CustomError) {
	const method = "Get"
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.l.Debug(
				method,
				"not found market in cache",
				"error", err.Error())
			return nil, nil
		}
		c.l.Error(layer, method, "failed to marshal data", err)
		return nil, errs.ErrUnavailableDataRedis
	}

	var result []dto.ViewMarketsResponse
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		c.l.Error(layer, method, "failed to unmarshal data", err)
		return nil, errs.ErrFailedDeserializeRedis
	}
	c.l.Debug(layer, method, "success get market cache")
	return result, nil
}

func (c *redisMarketsCache) Del(ctx context.Context, key string) *error.CustomError {
	const method = "Del"
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.l.Error(layer, method, "failed to delete market in cache", err)
		return errs.ErrDeleteRedis
	}
	c.l.Debug(layer, method, "success delete market in cache")
	return nil
}
