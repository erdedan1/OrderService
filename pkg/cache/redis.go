package cache

import (
	"OrderService/config"
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}

func NewRedisClient(config *config.Config) RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:         config.Infrastructure.RedisConfig.Host + ":" + config.Infrastructure.RedisConfig.Port,
		MinIdleConns: config.Infrastructure.RedisConfig.MinIdleConns,
		PoolSize:     config.Infrastructure.RedisConfig.PoolSize,
		PoolTimeout:  time.Duration(config.Infrastructure.RedisConfig.PoolTimeout),
		Password:     config.Infrastructure.RedisConfig.Password,
		DB:           config.Infrastructure.RedisConfig.DB,
	})

	return client
}
