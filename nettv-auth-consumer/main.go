package main

import (
	"context"

	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/api"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/db"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/logger"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/nettv-auth-consumer/services"
)

func main() {
	ctx := context.Background()

	NewConfig := config.NewConfig(logger.Logger())
	C := NewConfig.GetNewConfig()
	var rabbitmqsvc services.RabbitMQService
	{
		mongodbclient := db.NewMongoClient(
			logger.Logger(),
			C.Mongo,
		)

		client := mongodbclient.MongoClientConnect()

		redisclient := db.NewRedisClient(
			logger.Logger(),
			C.Redis,
		)

		rclient := redisclient.RedisConnect()

		mongoService := services.NewMongoDBService(
			logger.Logger(),
			C.Mongo,
			client,
		)

		redisService := services.NewRedisService(
			logger.Logger(),
			rclient,
		)

		serviceWlink := services.WlinkService(logger.Logger())

		nettvService := services.NewNettvService(logger.Logger())

		corporateService := api.CorporateService(logger.Logger())

		corporateNettvService := api.CorporateNettvService(logger.Logger())

		service := services.NewService(
			logger.Logger(),
			mongoService,
			redisService,
			serviceWlink,
			nettvService,
			corporateService,
			corporateNettvService,
		)

		rabbitmqsvc = services.NewRabbitMQService(
			logger.Logger(),
			service,
			C.RabbitMQ,
			C.RabbitMQEbill,
		)
	}

	rabbitmqsvc.QueueListen(ctx)
}
