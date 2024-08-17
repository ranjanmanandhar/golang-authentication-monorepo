package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
)

type RabbitMQService interface {
	QueueListen(ctx context.Context)
	FailOnError(ctx context.Context, err error, msg string)
}

type rabbitmqservice struct {
	logger   log.Logger
	service  service
	config   config.RabbitMQ
	configEb config.RabbitMQEbill
}

func NewRabbitMQService(logger log.Logger, service service, config config.RabbitMQ, configMq config.RabbitMQEbill) RabbitMQService {
	return rabbitmqservice{
		logger:   logger,
		service:  service,
		config:   config,
		configEb: configMq,
	}
}

func (r rabbitmqservice) FailOnError(ctx context.Context, err error, msg string) {
	logger := log.With(r.logger, "method", "QueueListen")
	if err != nil {
		fmt.Println(err)
		level.Error(logger).Log("err", err)
	}
}

func (r *rabbitmqservice) WaitForReconnect() {
	logger := log.With(r.logger, "method", "WaitForReconnect")

	delay := r.config.ReconnectDelay

	if delay.String() == "0s" {
		delay = time.Second * 30
	}

	level.Info(logger).Log("msg", fmt.Sprintf("Restarting connection in %s", delay.String()))

	time.Sleep(delay)
}

func (r rabbitmqservice) QueueListen(ctx context.Context) {
	for {
		rabbitmqConnectString := fmt.Sprintf("amqp://%s:%s@%s:%s/", r.config.Username, r.config.Password, r.config.HostName, r.config.Port)
		ebillRabbitmqConnectString := fmt.Sprintf("amqp://%s:%s@%s:%s/", r.configEb.Username, r.configEb.Password, r.configEb.HostName, r.configEb.Port)

		conn, err := amqp.Dial(rabbitmqConnectString)
		if err != nil {
			r.FailOnError(ctx, err, "Failed to connect to RabbitMQ")
			r.WaitForReconnect()
			continue
		}

		connEbillMq, err := amqp.Dial(ebillRabbitmqConnectString)
		if err != nil {
			r.FailOnError(ctx, err, "Failed to connect to Ebill RabbitMQ")
			r.WaitForReconnect()
			continue
		}

		notifyConnectionClose := conn.NotifyClose(make(chan *amqp.Error))

		// ebillNotifyConnectionClose := connEbillMq.NotifyClose(make(chan *amqp.Error))

		ch, err := conn.Channel()
		if err != nil {
			r.FailOnError(ctx, err, "Failed to open a channel")
			r.WaitForReconnect()
			continue
		}

		chEbill, err := connEbillMq.Channel()
		if err != nil {
			r.FailOnError(ctx, err, "Failed to open a channel in ebill rabbitmq")
			r.WaitForReconnect()
			continue
		}

		err = ch.ExchangeDeclare(
			r.config.ExchangeName, // name
			"fanout",              // type
			true,                  // durable
			false,                 // auto-deleted
			false,                 // internal
			false,                 // no-wait
			nil,                   // arguments
		)
		r.FailOnError(ctx, err, "Failed to declare an exchange")

		err = chEbill.ExchangeDeclare(
			r.configEb.ExchangeName, // name
			"fanout",                // type
			true,                    // durable
			false,                   // auto-deleted
			false,                   // internal
			false,                   // no-wait
			nil,                     // arguments
		)
		r.FailOnError(ctx, err, "Failed to declare an ebill exchange")

		q, err := ch.QueueDeclare(
			r.config.QueueName, // name
			false,              // durable
			false,              // delete when unused
			false,              // exclusive
			false,              // no-wait
			nil,                // arguments
		)
		r.FailOnError(ctx, err, "Failed to declare a queue")

		qEb, err := chEbill.QueueDeclare(
			r.configEb.QueueName, // name
			false,                // durable
			false,                // delete when unused
			false,                // exclusive
			false,                // no-wait
			nil,                  // arguments
		)
		r.FailOnError(ctx, err, "Failed to declare a ebill queue")

		err = ch.QueueBind(
			q.Name,                // queue name
			"",                    // routing key
			r.config.ExchangeName, // exchange
			false,
			nil,
		)
		r.FailOnError(ctx, err, "Failed to bind a queue")

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		r.FailOnError(ctx, err, "Failed to register a consumer")

		msgsEb, err := chEbill.Consume(
			qEb.Name, // queue
			"",       // consumer
			true,     // auto-ack
			false,    // exclusive
			false,    // no-local
			false,    // no-wait
			nil,      // args
		)
		r.FailOnError(ctx, err, "Failed to register a ebill consumer")

		forever := make(chan bool)

		go func() {
		P:
			for {
				select {
				case connErr := <-notifyConnectionClose:
					level.Error(r.logger).Log("msg", connErr)
					close(forever)
					break P
				case q1msgs := <-msgs:
					var newm map[string]interface{}
					err := json.Unmarshal(q1msgs.Body, &newm)
					if err != nil {
						fmt.Println(err)
					}
					r.service.ProcessQueue(ctx, q1msgs)
				case ebMsgs := <-msgsEb:
					var ebm map[string]interface{}
					err := json.Unmarshal(ebMsgs.Body, &ebm)
					if err != nil {
						fmt.Println(err)
					}
					r.service.ProcessQueue(ctx, ebMsgs)
				}
			}
		}()

		level.Info(r.logger).Log("msg", " [*] Waiting for messages. To exit press CTRL+C")
		<-forever
		r.WaitForReconnect()
	}
}
