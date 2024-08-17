package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/go-redis/redis"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
	libtype "gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type authRepository struct {
	redis  *redis.Client
	mongo  *mongo.Client
	config config.Mongo
	logger log.Logger
}

type AuthRepository interface {
	GetUser(username string) (libtype.Data, error)
	GetCorporateUser(username string) (libtype.CorporateData, error)
	GetCustomerCode(customerCode string) (types.CorporateCustomerCodeData, error)
	StoreRetailCustomer(payload types.Data) (bool, error)
	StoreCorporateCustomerCode(payload types.CorporateCustomerCodeData) (bool, error)
	StoreCorporateCustomerFallback(payload types.CorporateData) (bool, error)
	SeachByMacAddress(mac string) (libtype.Data, error)
	SearchByUsernameAndMac(username string, mac string) (libtype.OldSysData, error)
	Search(searchFilter bson.D) (libtype.OldSysData, error)
	SyncToken(Paylaod types.OldSysUser)
	SyncTokenInBatch(Payload []interface{})
	DeleteCustomer(filter bson.D)
	GetToken(searchFilter bson.D) (libtype.OldSysUser, error)
}

func New(r *redis.Client, m *mongo.Client, c config.Mongo, logger log.Logger) AuthRepository {
	return authRepository{
		redis:  r,
		mongo:  m,
		config: c,
		logger: logger,
	}
}

type CheckStatus struct {
	Authenticate bool `bson:"authenticate"`
	Status       int  `bson:"status"`
}

func (a authRepository) GetUser(username string) (libtype.Data, error) {
	coll := a.mongo.Database(a.config.DatabaseName).Collection(a.config.CollectionName)
	var result libtype.Data
	redisData, err := a.GetRetailUserInRedis(username)
	if err == nil {
		level.Error(a.logger).Log("METHOD", "GetUser", "msg", err)
		return redisData, nil
	}
	err = coll.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&result)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "GetUser", "msg", err)

		if err == mongo.ErrNoDocuments {
			return libtype.Data{}, err
		}
		panic(err)
	}
	a.SetRetailUserInRedis(result)
	return result, nil
}

func (a authRepository) GetCorporateUser(username string) (libtype.CorporateData, error) {
	coll := a.mongo.Database(a.config.DatabaseName).Collection(a.config.CollectionName)
	var result libtype.CorporateData
	a.GetCorporateUserInRedis(username)
	err := coll.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&result)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "GetCorporateUser", "msg", err)

		if err == mongo.ErrNoDocuments {
			return libtype.CorporateData{}, err
		}
		panic(err)
	}
	a.SetCorporateUserInRedis(result)
	return result, nil
}

func (a authRepository) GetCustomerCode(customerCode string) (types.CorporateCustomerCodeData, error) {
	coll := a.mongo.Database(a.config.DatabaseName).Collection("customer_codes")
	var result types.CorporateCustomerCodeData
	err := coll.FindOne(context.TODO(), bson.D{{"customer_code", customerCode}}).Decode(&result)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "GetCustomerCode", "msg", err)
		if err == mongo.ErrNoDocuments {
			return types.CorporateCustomerCodeData{}, err
		}
	}
	return result, nil
}

func (a authRepository) StoreRetailCustomer(payload types.Data) (bool, error) {
	collection := a.mongo.Database(a.config.DatabaseName).Collection(a.config.CollectionName)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"username": payload.Username,
	}
	update := bson.D{{"$set", payload}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "StoreRetailCustomer", "msg", err)
		return false, err
	}
	return true, nil
}

func (a authRepository) StoreCorporateCustomerCode(payload types.CorporateCustomerCodeData) (bool, error) {
	collection := a.mongo.Database(a.config.DatabaseName).Collection("customer_codes")
	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"customer_code": payload.CustomerCode,
	}
	update := bson.D{{"$set", payload}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "StoreCorporateCustomerCode", "msg", err)
		return false, err
	}
	return true, nil
}

func (a authRepository) StoreCorporateCustomerFallback(payload types.CorporateData) (bool, error) {
	collection := a.mongo.Database(a.config.DatabaseName).Collection(a.config.CollectionName)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"username": payload.Username,
	}
	update := bson.D{{"$set", payload}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "StoreCorporateCustomerFallback", "msg", err)
		return false, err
	}
	return true, nil
}

func (a authRepository) SetRetailUserInRedis(customerdata libtype.Data) {
	data, err := json.Marshal(customerdata)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "SetRetailUserInRedis", "msg", err)
	}
	a.redis.Set(customerdata.Username, data, time.Duration(1)*time.Second)

}

func (a authRepository) SetCorporateUserInRedis(customerdata libtype.CorporateData) {
	a.redis.Set(customerdata.Username, customerdata, time.Duration(60)*time.Second)
}

func (a authRepository) GetRetailUserInRedis(username string) (libtype.Data, error) {
	get, err := a.redis.Get(username).Bytes()
	var data libtype.Data
	if err != nil {
		level.Error(a.logger).Log("METHOD", "SetCorporateUserInRedis", "msg", err)
		return data, err
	}
	json.Unmarshal(get, &data)
	return data, nil
}

func (a authRepository) GetCorporateUserInRedis(username string) (libtype.CorporateData, error) {
	get, err := a.redis.Get(username).Bytes()
	var data libtype.CorporateData
	if err != nil {
		level.Error(a.logger).Log("METHOD", "GetCorporateUserInRedis", "msg", err)
		return data, err
	}
	json.Unmarshal(get, &data)
	return data, nil
}

func (a authRepository) SeachByMacAddress(mac string) (libtype.Data, error) {
	coll := a.mongo.Database(a.config.DatabaseName).Collection(a.config.CollectionName)
	var result types.Data
	err := coll.FindOne(context.TODO(), bson.D{{"nettvsettopboxes", mac}}).Decode(&result)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "SeachByMacAddress", "msg", err)
		return result, err
	}
	return result, nil
}

func (a authRepository) SearchByUsernameAndMac(username string, mac string) (libtype.OldSysData, error) {
	coll := a.mongo.Database(a.config.DatabaseName).Collection(a.config.CollectionName)
	var result types.OldSysData
	err := coll.FindOne(context.TODO(), bson.D{{"nettvsettopboxes", mac}, {"username", username}}).Decode(&result)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "SearchByUsernameAndMac", "msg", err)
		return result, err
	}
	return result, nil
}

func (a authRepository) Search(searchFilter bson.D) (libtype.OldSysData, error) {
	coll := a.mongo.Database(a.config.DatabaseName).Collection(a.config.CollectionName)
	var result types.OldSysData
	err := coll.FindOne(context.TODO(), searchFilter).Decode(&result)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "Search", "msg", err)
		return result, err
	}
	return result, nil
}

func (a authRepository) SyncToken(Paylaod types.OldSysUser) {
	collection := a.mongo.Database(a.config.DatabaseName).Collection("old-system-token")
	opts := options.Update().SetUpsert(true)
	filter := bson.M{
		"username":   Paylaod.Username,
		"stb_box_id": Paylaod.StbBoxId,
	}
	update := bson.D{{"$set", Paylaod}}
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "SyncToken", "msg", err)
		// return false, err
	}
	// return true, nil
}

func (a authRepository) SyncTokenInBatch(Payload []interface{}) {
	collection := a.mongo.Database(a.config.DatabaseName).Collection("old-system-token")
	_, err := collection.InsertMany(context.Background(), Payload)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "SyncTokenInBatch", "msg", err)
	}
}

func (a authRepository) GetToken(searchFilter bson.D) (libtype.OldSysUser, error) {
	fmt.Println("get token", searchFilter)
	coll := a.mongo.Database(a.config.DatabaseName).Collection("old-system-token")
	var result libtype.OldSysUser
	err := coll.FindOne(context.TODO(), searchFilter).Decode(&result)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "GetToken", "msg", err)

		if err == mongo.ErrNoDocuments {
			return libtype.OldSysUser{}, err
		}
		panic(err)
	}
	return result, nil
}

func (a authRepository) DeleteCustomer(filter bson.D) {
	collection := a.mongo.Database(a.config.DatabaseName).Collection(a.config.CollectionName)
	result, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "DeleteCustomer", "msg", err)
	}
	fmt.Printf("Number of documents deleted: %d\n", result.DeletedCount)
}
