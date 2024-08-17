package services

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-redis/redis"
)

type RedisService interface {
	DelFromRedis(ctx context.Context, username string)
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

func (r redisservice) DelFromRedis(ctx context.Context, username string) {
	r.rclient.Del(username)
	level.Info(r.logger).Log("username", username)

}
