package config

// import (
// 	"flag"
// 	"time"

// 	"github.com/spf13/viper"
// 	_ "github.com/spf13/viper/remote"
// )

// var RuntimeViper = viper.New()
// var consulAddr = flag.String("consulAddress", "localhost:8500", "Consul URL")
// var consulKey = flag.String("consulKey", "nettv-auth", "Consul Key")

// type RabbitMQ struct {
// 	HostName       string        `json:"hostname" mapstructure:"hostname"`
// 	Port           string        `json:"port" mapstructure:"port"`
// 	QueueName      string        `json:"queueName" mapstructure:"queueName"`
// 	Username       string        `json:"username" mapstructure:"username"`
// 	Password       string        `json:"password" mapstructure:"password"`
// 	ExchangeName   string        `json:"exchangeName" mapstructure:"exchangeName"`
// 	ReconnectDelay time.Duration `json:"reconnectDelay" mapstructure:"reconnectDelay"`
// 	UseAuth        bool          `json:"useAuth"`
// }

// type Mongo struct {
// 	HostName       string `json:"hostname" mapstructure:"hostname"`
// 	Port           string `json:"port" mapstructure:"port"`
// 	Username       string `json:"username" mapstructure:"username"`
// 	Password       string `json:"password" mapstructure:"password"`
// 	DatabaseName   string `json:"database_name" mapstructure:"database_name"`
// 	CollectionName string `json:"collection_name" mapstructure:"collection_name"`
// 	UseAuth        bool   `json:"useAuth"`
// }

// type Redis struct {
// 	HostName string `json:"hostname" mapstructure:"hostname"`
// 	Port     string `json:"port" mapstructure:"port"`
// 	UseAuth  bool   `json:"useAuth"`
// 	Username string `json:"username" mapstructure:"username"`
// 	Password string `json:"password" mapstructure:"password"`
// }

// type Config struct {
// 	RabbitMQ RabbitMQ `json:"rabbitmq" mapstructure:"rabbitmq"`
// 	Mongo    Mongo    `json:"mongo" mapstructure:"mongo"`
// 	Redis    Redis    `json:"redis" mapstructure:"redis"`
// }

// func GetNewConfig() Config {
// 	var ConfigData Config
// 	RuntimeViper.AddRemoteProvider("consul", *consulAddr, *consulKey)
// 	RuntimeViper.SetConfigType("json")
// 	err := RuntimeViper.ReadRemoteConfig()

// 	if err != nil {
// 		panic(err)
// 	}

// 	err = RuntimeViper.Unmarshal(&ConfigData)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return ConfigData
// }
