package service

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"time"

// 	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // type MongoClient interface {
// // 	MongoClientConnect() *mongo.Client
// // 	Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc)
// // }

// // type mongoClient struct {
// // 	config config.Mongo
// // }

// // func NewMongoClient(config config.Mongo) MongoClient {
// // 	return mongoClient{
// // 		config: config,
// // 	}
// // }

// // func (m mongoClient) MongoClientConnect() *mongo.Client {
// // 	connectString := fmt.Sprintf("mongodb://%s:%s", m.config.HostName, m.config.Port)
// // 	if m.config.UseAuth {
// // 		connectString = fmt.Sprintf("mongodb://%s:%s@%s:%s", m.config.Username, m.config.Password, m.config.HostName, m.config.Port)
// // 	}

// // 	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
// // 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectString))

// // 	if err != nil {
// // 		fmt.Println(err)
// // 		// level.Error(Newlogger).Log(constant.LOG_STATUS, err)
// // 		os.Exit(-1)
// // 	}
// // 	fmt.Println("database connection successful")
// // 	// level.Info(Newlogger).Log(constant.LOG_STATUS, "connect to db success")

// // 	return client
// // }

// // func (m mongoClient) Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

// // 	defer cancel()

// // 	defer func() {
// // 		if err := client.Disconnect(ctx); err != nil {
// // 			panic(err)
// // 		}
// // 	}()

// // }
