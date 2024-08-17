package queue

import (
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/streadway/amqp"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
)

type RabbitmqClient interface {
	HealthCheck() types.RedisHealthCheckResponse
	HealthCheckEbillRabbitmq() types.RedisHealthCheckResponse
}

type rabbitmqClient struct {
	logger log.Logger
	config config.RabbitMQ
}

func NewRabbitmqClient(logger log.Logger, config config.RabbitMQ) RabbitmqClient {
	return rabbitmqClient{
		logger: logger,
		config: config,
	}
}

func (r rabbitmqClient) HealthCheck() types.RedisHealthCheckResponse {
	var healthCheck types.RedisHealthCheckResponse
	var checkErr error
	rabbitmqConnectString := fmt.Sprintf("amqp://%s:%s@%s:%s/", r.config.Username, r.config.Password, r.config.HostName, r.config.Port)
	conn, err := amqp.DialConfig(rabbitmqConnectString, amqp.Config{
		Dial: amqp.DefaultDial(200 * time.Millisecond),
	})
	if err != nil {
		checkErr = fmt.Errorf("RabbitMQ health check failed on dial phase: %w", err)
		healthCheck.Status = "failed"
		healthCheck.Message = "rabbimq connection failed"
		return healthCheck
	}
	defer func() {
		// override checkErr only if there were no other errors
		if err := conn.Close(); err != nil && checkErr == nil {
			checkErr = fmt.Errorf("RabbitMQ health check failed to close connection: %w", err)
		}
	}()
	healthCheck.Status = "success"
	healthCheck.Message = "rabbimq connection success"
	return healthCheck
}

func (r rabbitmqClient) HealthCheckEbillRabbitmq() types.RedisHealthCheckResponse {
	var healthCheck types.RedisHealthCheckResponse
	return healthCheck
}
