package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	LOG_TYPE             = "log_type"
	APPLICATION_LOG_TYPE = "application"
	LOG_SOURCE           = "source"
	MONGO_LOG_SOURCE     = "mongo"
	REDIS_LOG_SOURCE     = "redis"
	LOG_STATUS           = "status"
	METHOD               = "method"
	ENV_REMOTE           = "remote"
)

type MongoClient interface {
	MongoClientConnect() *mongo.Client
	Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc)
	MongoHealthCheck(ctx context.Context) types.RedisHealthCheckResponse
	NewMongoHealthCheckConnection() error
}

type mongoClient struct {
	logger log.Logger
	config config.Mongo
}

func NewMongoClient(logger log.Logger, config config.Mongo) MongoClient {
	return mongoClient{
		logger: logger,
		config: config,
	}
}

func (m mongoClient) MongoClientConnect() *mongo.Client {
	Newlogger := log.With(m.logger,
		LOG_TYPE, APPLICATION_LOG_TYPE,
		LOG_SOURCE, MONGO_LOG_SOURCE,
		METHOD, "MongoClientConnect",
	)
	connectString := fmt.Sprintf("mongodb://%s:%s", m.config.HostName, m.config.Port)
	if m.config.Environment == "production" {
		connectString = fmt.Sprintf("mongodb://%s:%s@%s:%s", m.config.Username, m.config.Password, m.config.HostName, m.config.Port)
	}
	level.Info(Newlogger).Log("Db connection string", connectString)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString))

	if err != nil {

		level.Error(Newlogger).Log(LOG_STATUS, err)
		os.Exit(-1)
	}

	level.Info(Newlogger).Log(LOG_STATUS, "connect to db success")

	return client
}

func (m mongoClient) Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

}

func (m mongoClient) MongoHealthCheck(ctx context.Context) types.RedisHealthCheckResponse {
	respErr := m.NewMongoHealthCheckConnection()
	var redisHealthcheckResponse types.RedisHealthCheckResponse

	if respErr != nil {
		level.Error(m.logger).Log(LOG_STATUS, respErr)
		redisHealthcheckResponse.Status = "failed"
		redisHealthcheckResponse.Message = "mongo connection failed"
		fmt.Println("resp", redisHealthcheckResponse)
		return redisHealthcheckResponse
	}

	redisHealthcheckResponse.Status = "success"
	redisHealthcheckResponse.Message = "mongo connection success"
	return redisHealthcheckResponse
}

func (m mongoClient) NewMongoHealthCheckConnection() error {
	connectString := fmt.Sprintf("mongodb://%s:%s", m.config.HostName, m.config.Port)
	if m.config.Environment == "production" {
		connectString = fmt.Sprintf("mongodb://%s:%s@%s:%s", m.config.Username, m.config.Password, m.config.HostName, m.config.Port)
	}
	var checkErr error
	TimeoutConnect := 10 * time.Second

	TimeoutDisconnect := 10 * time.Second

	TimeoutPing := 10 * time.Second

	client, err := mongo.NewClient(options.Client().ApplyURI(connectString))
	if err != nil {
		fmt.Println("mongoDB health check failed on client creation:")
		checkErr := fmt.Errorf("mongoDB health check failed on client creation: %w", err)
		return checkErr
	}

	ctxConn, cancelConn := context.WithTimeout(context.Background(), TimeoutConnect)
	defer cancelConn()

	err = client.Connect(ctxConn)
	if err != nil {
		fmt.Println("mongoDB health check failed on connect")
		checkErr := fmt.Errorf("mongoDB health check failed on connect: %w", err)
		return checkErr
	}

	defer func() {
		ctxDisc, cancelDisc := context.WithTimeout(context.Background(), TimeoutDisconnect)
		defer cancelDisc()
		// override checkErr only if there were no other errors
		if err := client.Disconnect(ctxDisc); err != nil && checkErr == nil {
			fmt.Println("mongoDB health check failed on closing connection:")
			checkErr = fmt.Errorf("mongoDB health check failed on closing connection: %w", err)
		}
	}()

	ctxPing, cancelPing := context.WithTimeout(context.Background(), TimeoutPing)
	defer cancelPing()

	err = client.Ping(ctxPing, readpref.Primary())
	if err != nil {
		checkErr = fmt.Errorf("mongoDB health check failed on pingss: %w", err)
		return checkErr
	}

	return nil
}
