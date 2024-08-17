package service

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-redis/redis/v8"
)

type RedisService interface {
	SetCache(ctx context.Context, username string, data map[string]interface{})
}

type redisservice struct {
	logger  log.Logger
	rclient *redis.Client
}

func NewRedisService(logger log.Logger, rclient *redis.Client) RedisService {
	return redisservice{
		logger:  logger,
		rclient: rclient,
	}
}

func (r redisservice) SetCache(ctx context.Context, username string, data map[string]interface{}) {
	err := r.rclient.Set(ctx, username, data, 0).Err()
	if err != nil {
		panic(err)
	}
}
