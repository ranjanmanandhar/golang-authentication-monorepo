package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log/level"
	_ "gitlab-server.wlink.com.np/nettv/nettv-auth/authenticator/cmd/app/docs"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/authenticator/internal/controller"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/authenticator/internal/repository"

	"gitlab-server.wlink.com.np/nettv/nettv-auth/authenticator/internal/service"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/api"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/config"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/db"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/logger"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/queue"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	logger := logger.Logger()
	config := config.NewConfig(logger)
	C := config.GetNewConfig()
	level.Info(logger).Log("redis")

	redis := db.NewRedisClient(logger, C.Redis)
	redisClient := redis.RedisConnect()

	mongoDb := db.NewMongoClient(logger, C.Mongo)
	mongoClient := mongoDb.MongoClientConnect()

	oracleConnect := db.NewOracleClient(logger, C.Oracle)
	oracleDB := oracleConnect.ConnectOracle()

	rabbitmqSrv := queue.NewRabbitmqClient(logger, C.RabbitMQ)

	authRepository := repository.New(redisClient, mongoClient, C.Mongo, logger)

	nettvservice := service.NewNettvService(logger, service.Nettv{
		Hostname: "www3.wlink.com.np",
		Schema:   "https",
		Token:    "randomstring",
	})

	customerLoginService := service.NewCustomerLogin(logger)

	apiWlink := api.ApiWlinkService(logger)

	authController := controller.New(authRepository, nettvservice, customerLoginService, logger)
	oldController := controller.NewOldSysAuth(authRepository, nettvservice, customerLoginService, apiWlink, logger, oracleDB)
	healthController := controller.NewHealthController(redis, mongoDb, rabbitmqSrv, C.Mongo, logger)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "welcome to authenticator service")
	})

	v1 := r.Group("/api")
	{
		client := v1.Group("client")
		{
			client.POST("/default_authenticate", Mymiddleware(C.Secret), authController.DefaultAuthenticate)
			client.POST("/verify_client", authController.VerifyClient)
			client.POST("/check_user_status", Mymiddleware(C.Secret), authController.CheckUserStatus)
			client.POST("/macwise_user_check", oldController.MacWiseUserCheck)
			client.POST("/adl_authentication", oldController.AdlAuthentication)
			client.POST("/old_sys_authentcation", oldController.Authentication)
			client.POST("/old_sys_token_sync", oldController.OldSysTokenSync)
		}
		v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	health := r.Group("/health")
	{
		health.GET("/", healthController.HealthCheckDependencies)
		health.GET("/mongo", healthController.CheckMongoHealth)
		health.GET("/redis", healthController.CheckRedisHealth)
		health.GET("/rabbitmq", healthController.CheckRabbitmqHealth(C))
	}

	return r
}

func Mymiddleware(config config.Secret) gin.HandlerFunc {
	logger := logger.Logger()
	return func(c *gin.Context) {
		level.Info(logger).Log("msg", "middleware here")
		Appsecret := c.Request.Header.Get("APPSECRET")
		Appid := c.Request.Header.Get("APPID")

		if Appsecret == "" && Appid == "" {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error_code": http.StatusUnprocessableEntity, "error_message": "Header Unprocessable entity"})
			return
		}

		if Appsecret != config.AppSecret || Appid != config.AppId {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": http.StatusBadRequest, "error_message": "Invalid App Headers"})
			return
		}

		c.Next()
	}
}
