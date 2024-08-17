package config

import (
	"flag"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/util"
)

var RuntimeViper = viper.New()
var consulAddr = flag.String("consulAddress", "localhost:8500", "Consul URL")
var consulKey = flag.String("consulKey", "nettv-auth", "Consul Key")
var envFlag = flag.Bool("env", false, "Environment Type")

type ConfigInterface interface {
	GetNewConfig() Config
}

type config struct {
	logger log.Logger
}

func NewConfig(logger log.Logger) ConfigInterface {
	return config{
		logger: logger,
	}
}

type RabbitMQ struct {
	HostName       string        `json:"hostname" mapstructure:"hostname"`
	Port           string        `json:"port" mapstructure:"port"`
	QueueName      string        `json:"queueName" mapstructure:"queueName"`
	Username       string        `json:"username" mapstructure:"username"`
	Password       string        `json:"password" mapstructure:"password"`
	ExchangeName   string        `json:"exchangeName" mapstructure:"exchangeName"`
	ReconnectDelay time.Duration `json:"reconnectDelay" mapstructure:"reconnectDelay"`
}

type RabbitMQEbill struct {
	HostName     string `json:"eb-hostname" mapstructure:"eb-hostname"`
	Port         string `json:"eb-port" mapstructure:"eb-port"`
	QueueName    string `json:"eb-queueName" mapstructure:"eb-queueName"`
	Username     string `json:"eb-username" mapstructure:"eb-username"`
	Password     string `json:"eb-password" mapstructure:"eb-password"`
	ExchangeName string `json:"eb-exchangeName" mapstructure:"eb-exchangeName"`
}

type Mongo struct {
	HostName       string `json:"hostname" mapstructure:"hostname"`
	Port           string `json:"port" mapstructure:"port"`
	Username       string `json:"username" mapstructure:"username"`
	Password       string `json:"password" mapstructure:"password"`
	DatabaseName   string `json:"database_name" mapstructure:"database_name"`
	CollectionName string `json:"collection_name" mapstructure:"collection_name"`
	Environment    string `json:"environment" mapstructure:"environment"`
}

type Redis struct {
	HostName string `json:"hostname" mapstructure:"hostname"`
	Port     string `json:"port" mapstructure:"port"`
}

type Oracle struct {
	HostName     string `json:"hostname" mapstructure:"hostname"`
	Port         string `json:"port" mapstructure:"port"`
	Username     string `json:"username" mapstructure:"username"`
	Password     string `json:"password" mapstructure:"password"`
	DatabaseName string `json:"database_name" mapstructure:"database_name"`
	ServiceName  string `json:"service_name" mapstructure:"service_name"`
	Environment  string `json:"environment" mapstructure:"environment"`
}

type Secret struct {
	AppId     string `json:"appid" mapstructure:"appid"`
	AppSecret string `json:"appsecret" mapstructure:"appsecret"`
}

type Config struct {
	RabbitMQ      RabbitMQ      `json:"rabbitmq" mapstructure:"rabbitmq"`
	RabbitMQEbill RabbitMQEbill `json:"rabbitmq-ebill" mapstructure:"rabbitmq-ebill"`
	Mongo         Mongo         `json:"mongo" mapstructure:"mongo"`
	Redis         Redis         `json:"redis" mapstructure:"redis"`
	Secret        Secret        `json:"secret" mapstructure:"secret"`
	Oracle        Oracle        `json:"oracle" mapstructure:"oracle"`
}

func (c config) GetNewConfig() Config {
	var ConfigData Config
	if !*envFlag {
		RuntimeViper.AddRemoteProvider("consul", *consulAddr, *consulKey)
		RuntimeViper.SetConfigType("json")
		err := RuntimeViper.ReadRemoteConfig()

		if err != nil {
			level.Error(c.logger).Log("err", err)
		}

		err = RuntimeViper.Unmarshal(&ConfigData)
		if err != nil {
			level.Error(c.logger).Log("err", err)
		}
	} else {
		ConfigData = Config{
			RabbitMQ{
				HostName:     util.GetEnvValue("RABBITMQ_HOSTNAME", "10.12.7.219"),
				Port:         util.GetEnvValue("RABBITMQ_PORT", "35672"),
				QueueName:    util.GetEnvValue("RABBITMQ_QUEUENAME", "nettv_auth_consumer"),
				Username:     util.GetEnvValue("RABBITMQ_USERNAME", "guest"),
				Password:     util.GetEnvValue("RABBITMQ_PASSWORD", "guest"),
				ExchangeName: util.GetEnvValue("RABBITMQ_EXCHANGENAME", "nettv_auth_consumer"),
			},
			RabbitMQEbill{
				HostName:     util.GetEnvValue("EB_RABBITMQ_HOSTNAME", "10.12.7.108"),
				Port:         util.GetEnvValue("EB_RABBITMQ_PORT", "5672"),
				QueueName:    util.GetEnvValue("EB_RABBITMQ_QUEUENAME", "nettv_auth_consumer"),
				Username:     util.GetEnvValue("EB_RABBITMQ_USERNAME", "guest"),
				Password:     util.GetEnvValue("EB_RABBITMQ_PASSWORD", "guest"),
				ExchangeName: util.GetEnvValue("EB_RABBITMQ_EXCHANGENAME", "nettv_auth_consumer"),
			},
			Mongo{
				HostName:       util.GetEnvValue("MONGO_HOSTNAME", "10.12.7.219"),
				Port:           util.GetEnvValue("MONGO_PORT", "27017"),
				Username:       util.GetEnvValue("MONGO_USERNAME", "admin"),
				Password:       util.GetEnvValue("MONGO_PASSWORD", "password"),
				DatabaseName:   util.GetEnvValue("MONGO_DBNAME", "nettv_auth_consumer"),
				CollectionName: util.GetEnvValue("MONGO_COLLECTIONNAME", "nettv_auth_consumer"),
				Environment:    util.GetEnvValue("MONGO_ENVIRONMENT", "development"),
			},
			Redis{
				HostName: util.GetEnvValue("REDIS_HOSTNAME", "10.12.7.219"),
				Port:     util.GetEnvValue("REDIS_PORT", "6381"),
			},
			Secret{
				AppId:     util.GetEnvValue("APP_ID", "NETTV"),
				AppSecret: util.GetEnvValue("APP_Secret", "e0fe1cc0f2dad4c68646b1f7bea87dc6"),
			},
			Oracle{
				HostName:     util.GetEnvValue("ORACLE_HOSTNAME", "raddbdev-scan.wlink.com.np"),
				Port:         util.GetEnvValue("ORACLE_PORT", "1521"),
				Username:     util.GetEnvValue("ORACLE_USERNAME", "ebill"),
				Password:     util.GetEnvValue("ORACLE_PASSWORD", "Orcl_4Dev"),
				DatabaseName: util.GetEnvValue("ORACLE_DBNAME", "raddb"),
				ServiceName:  util.GetEnvValue("ORACLE_SERVICENAME", "raddb"),
				Environment:  util.GetEnvValue("ORACLE_ENVIRONMENT", "development"),
			},
		}
	}
	return ConfigData
}
