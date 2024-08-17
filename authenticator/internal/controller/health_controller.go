package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/streadway/amqp"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/db"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/queue"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
)

type healthController struct {
	redis    db.RedisCLient
	mongo    db.MongoClient
	rabbitmq queue.RabbitmqClient
	config   config.Mongo
	logger   log.Logger
}

type HealthController interface {
	CheckRedisHealth(c *gin.Context)
	CheckMongoHealth(c *gin.Context)
	HealthCheckDependencies(c *gin.Context)
	CheckRabbitmqHealth(config config.Config) gin.HandlerFunc
}

func NewHealthController(r db.RedisCLient, m db.MongoClient, rMq queue.RabbitmqClient, c config.Mongo, l log.Logger) HealthController {
	return healthController{
		redis:    r,
		mongo:    m,
		rabbitmq: rMq,
		config:   c,
		logger:   l,
	}
}

func (h healthController) CheckRedisHealth(c *gin.Context) {
	resp := h.redis.RedisHealthCheck()
	c.JSON(http.StatusOK, resp)
}

func (h healthController) CheckRabbitmqHealth(config config.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var healthCheck types.RedisHealthCheckResponse
		var checkErr error
		rabbitmqConnectString := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitMQ.Username, config.RabbitMQ.Password, config.RabbitMQ.HostName, config.RabbitMQ.Port)
		conn, err := amqp.DialConfig(rabbitmqConnectString, amqp.Config{
			Dial: amqp.DefaultDial(200 * time.Millisecond),
		})
		if err != nil {
			checkErr = fmt.Errorf("RabbitMQ health check failed on dial phase: %w", err)
			healthCheck.Status = "failed"
			healthCheck.Message = "rabbimq connection failed"
			c.JSON(http.StatusOK, healthCheck)
			return
		}
		defer func() {
			// override checkErr only if there were no other errors
			if err := conn.Close(); err != nil && checkErr == nil {
				checkErr = fmt.Errorf("RabbitMQ health check failed to close connection: %w", err)
			}
		}()
		healthCheck.Status = "success"
		healthCheck.Message = "rabbimq connection success"
		c.JSON(http.StatusOK, healthCheck)
		return
	}
	return gin.HandlerFunc(fn)
}

func (h healthController) CheckMongoHealth(c *gin.Context) {
	resp := h.mongo.MongoHealthCheck(c)
	fmt.Println("resp here", resp)
	c.JSON(http.StatusAccepted, resp)
	return
}

func (h healthController) HealthCheckDependencies(c *gin.Context) {
	redisRes := h.redis.RedisHealthCheck()
	rabbitmqResp := h.rabbitmq.HealthCheck()
	mongoResp := h.mongo.MongoHealthCheck(c)
	c.JSON(http.StatusOK, gin.H{"redis": redisRes, "rabbitmq": rabbitmqResp, "mongo": mongoResp})
}
