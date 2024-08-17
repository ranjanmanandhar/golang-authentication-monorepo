package services

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBService interface {
	UpdateOrInsertCustomer(ctx context.Context, filter bson.M, update types.Data)
	UpdateOrInsertCustomerInfo(ctx context.Context, filter bson.M, update types.CustomerInfoData)
	UpdateOrInsertNettv(ctx context.Context, filter bson.M, update types.RetailNettvData)
	UpdateOrInsertCorporate(ctx context.Context, filter bson.M, update types.CircuitData)
	UpdateOrInsertCorporateCodes(ctx context.Context, filter bson.M, update types.CorporateCustomerCodeData)
	UpdateOrInsertCorporateNettv(ctx context.Context, filter bson.M, update types.CorporateNettvData)
	FindOne(ctx context.Context, filter bson.M) (bson.M, error)
}

type mongodbservice struct {
	logger log.Logger
	config config.Mongo
	client *mongo.Client
}

func NewMongoDBService(logger log.Logger, config config.Mongo, client *mongo.Client) MongoDBService {
	return mongodbservice{
		logger: logger,
		config: config,
		client: client,
	}
}

func (m mongodbservice) UpdateOrInsertCustomer(ctx context.Context, filter bson.M, updateData types.Data) {
	collection := m.client.Database(m.config.DatabaseName).Collection(m.config.CollectionName)
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", updateData}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}
}

func (m mongodbservice) UpdateOrInsertCustomerInfo(ctx context.Context, filter bson.M, updateData types.CustomerInfoData) {
	collection := m.client.Database(m.config.DatabaseName).Collection(m.config.CollectionName)
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", updateData}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}
}

func (m mongodbservice) UpdateOrInsertCorporate(ctx context.Context, filter bson.M, updateData types.CircuitData) {
	collection := m.client.Database(m.config.DatabaseName).Collection(m.config.CollectionName)
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", updateData}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}
}

func (m mongodbservice) UpdateOrInsertCorporateNettv(ctx context.Context, filter bson.M, updateData types.CorporateNettvData) {
	collection := m.client.Database(m.config.DatabaseName).Collection(m.config.CollectionName)
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", updateData}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}
}

func (m mongodbservice) UpdateOrInsertCorporateCodes(ctx context.Context, filter bson.M, updateData types.CorporateCustomerCodeData) {
	collection := m.client.Database(m.config.DatabaseName).Collection("customer_codes")
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", updateData}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}
}

func (m mongodbservice) UpdateOrInsertNettv(ctx context.Context, filter bson.M, updateData types.RetailNettvData) {
	collection := m.client.Database(m.config.DatabaseName).Collection(m.config.CollectionName)
	opts := options.Update().SetUpsert(true)
	update := bson.D{{"$set", updateData}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		panic(err)
	}
}

func (m mongodbservice) FindOne(ctx context.Context, filter bson.M) (bson.M, error) {
	collection := m.client.Database(m.config.DatabaseName).Collection(m.config.CollectionName)
	var result bson.M
	if err := collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			level.Error(m.logger).Log("msg", err)
			return nil, err
		}
		panic(err)
	}
	return result, nil
}
