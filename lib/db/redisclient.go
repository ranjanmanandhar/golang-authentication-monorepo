package db

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-redis/redis"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
)

type RedisCLient interface {
	RedisConnect() *redis.Client
	RedisHealthCheck() types.RedisHealthCheckResponse
}

type redisClient struct {
	logger log.Logger
	config config.Redis
}

func NewRedisClient(logger log.Logger, config config.Redis) RedisCLient {
	return redisClient{
		logger: logger,
		config: config,
	}
}

func (r redisClient) RedisConnect() *redis.Client {
	Newlogger := log.With(r.logger,
		LOG_TYPE, APPLICATION_LOG_TYPE,
		LOG_SOURCE, REDIS_LOG_SOURCE,
		METHOD, "RedisConnect",
	)

	redisConnectionString := fmt.Sprintf("%s:%s", r.config.HostName, r.config.Port)

	redisclient := redis.NewClient(&redis.Options{
		Addr: redisConnectionString,
	})

	pong, err := redisclient.Ping().Result()

	if err != nil {
		level.Error(Newlogger).Log(LOG_STATUS, err)
		// os.Exit(-1)
	}

	level.Info(Newlogger).Log(LOG_STATUS, pong)

	return redisclient
}

func (r redisClient) RedisHealthCheck() types.RedisHealthCheckResponse {
	Newlogger := log.With(r.logger,
		LOG_TYPE, APPLICATION_LOG_TYPE,
		LOG_SOURCE, REDIS_LOG_SOURCE,
		METHOD, "RedisHealthCheck",
	)
	var redisHealthcheckResponse types.RedisHealthCheckResponse
	redisConnectionString := fmt.Sprintf("%s:%s", r.config.HostName, r.config.Port)

	redisclient := redis.NewClient(&redis.Options{
		Addr: redisConnectionString,
	})

	pong, err := redisclient.Ping().Result()

	if err != nil {
		level.Error(Newlogger).Log(LOG_STATUS, err)
		redisHealthcheckResponse.Status = "failed"
		redisHealthcheckResponse.Message = "redis connection failed"
		return redisHealthcheckResponse
	}

	level.Info(Newlogger).Log(LOG_STATUS, pong)
	redisHealthcheckResponse.Status = "success"
	redisHealthcheckResponse.Message = "redis connection success"
	return redisHealthcheckResponse
}
