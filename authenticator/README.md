Authenticator service 


Cousul - nettv-auth

{
 "rabbitmq": {
        "hostname": "localhost",
        "port": "5671",
        "username": "guest",
        "password": "guest",
        "exchangeName": "nettv_auth_consumer",
        "queueName": "nettv_auth_consumer"
    },
   "rabbitmq-ebill": {
        "eb-hostname": "10.12.7.108",
        "eb-port": "5672",
        "eb-username": "guest",
        "eb-password": "guest",
        "eb-exchangeName": "ebill",
        "eb-queueName": "nettv_auth_consumer"
    },
  "mongo": {
 				"hostname" : "10.12.7.219",
       	"port" : "27017",
      	"username" : "admin",
      	"password" : "password",
       	"database_name" : "nettv_auth_consumer",
       	"collection_name" : "nettv_auth_consumer"
  	},
  "redis": {
    		"hostname" : "10.12.7.219",
    		"port": "36379"
  	},
  "secret":{
    "appid":"NETTV",
    "appsecret":"e0fe1cc0f2dad4c68646b1f7bea87dc6"
  }
}