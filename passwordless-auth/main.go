package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-kit/log/level"
	_ "github.com/godror/godror"
	"github.com/spf13/viper"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/logger"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/passwordless-auth/config"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/passwordless-auth/services"
)

const defaultPort string = "9100"

// func init() {
// 	logger := logger.Logger()
// 	logger = log.NewLogfmtLogger(os.Stderr)
// 	logger = log.NewSyncLogger(logger.Logger())
// 	logger = log.With(logger,
// 		"service", "Nettv Passwordless Authentication",
// 		"time", log.DefaultTimestamp,
// 		"caller", log.DefaultCaller,
// 	)
// }

func main() {
	var envFile string
	var dbConfig config.DBConfig

	flag.StringVar(&envFile, "envFile", "", "Use operating system environment variables")
	flag.Parse()

	if len(envFile) > 0 {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			level.Error(logger.Logger()).Log("Error Reading config file", err)
			os.Exit(1)
		}
	} else {
		viper.SetDefault("DB_CONNECTION", "")
		viper.SetDefault("DB_HOST", "")
		viper.SetDefault("DB_PORT", "")
		viper.SetDefault("DB_DATABASE", "")
		viper.SetDefault("DB_USERNAME", "")
		viper.SetDefault("DB_PASSWORD", "")
		viper.SetDefault("ESERVICE_DB_CONNECTION", "")
		viper.SetDefault("ESERVICE_DB_HOST", "")
		viper.SetDefault("ESERVICE_DB_PORT", "")
		viper.SetDefault("ESERVICE_DB_DATABASE", "")
		viper.SetDefault("ESERVICE_DB_USERNAME", "")
		viper.SetDefault("ESERVICE_DB_PASSWORD", "")
		viper.SetDefault("EBILL_DB_CONNECTION", "")
		viper.SetDefault("EBILL_DB_HOST", "raddbdev-scan.wlink.com.np")
		viper.SetDefault("EBILL_DB_PORT", "1521")
		viper.SetDefault("EBILL_DB_DATABASE", "raddb")
		viper.SetDefault("EBILL_DB_USERNAME", "ebill")
		viper.SetDefault("EBILL_DB_PASSWORD", "Orcl_4Dev")
		viper.AutomaticEnv()
	}

	if err := viper.Unmarshal(&dbConfig); err != nil {
		level.Error(logger.Logger()).Log("Env Variable Unmarshall Error", err)
		os.Exit(1)
	}

	level.Info(logger.Logger()).Log("Config", fmt.Sprintf("%+v", dbConfig))

	eservice, eserviceError := services.Connection("eservice", &dbConfig)

	if eserviceError != nil {
		level.Error(logger.Logger()).Log("DB Connection Error", eserviceError)
	}

	err := eservice.Ping()
	if err != nil {
		level.Error(logger.Logger()).Log("DB Ping Error", err)
	}

	ebill, ebillError := services.Connection("ebill", &dbConfig)

	if ebillError != nil {
		level.Error(logger.Logger()).Log("Ebill DB Connection Error", ebillError)
	}

	err = ebill.Ping()
	if err != nil {
		level.Error(logger.Logger()).Log("Ebill DB Ping Error", err)
	}

	auth := services.AuthService(eservice, ebill, logger.Logger())

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/match_ip_mac", func(w http.ResponseWriter, r *http.Request) {

		mac := strings.ToUpper(r.URL.Query().Get("mac"))
		stb_ip_address := r.URL.Query().Get("ip_address")
		if stb_ip_address == "" {
			stb_ip_address = services.GetIP(r)
		}
		secret, id := services.FetchToken(r)

		if secret == "" || id == "" {
			level.Error(logger.Logger()).Log("Error", "App Secret and AppID both must be present")
			data := map[string]interface{}{
				"error_message": "Header Unprocessable entity",
				"error_code":    422,
			}
			w.WriteHeader(422)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
			return
		}
		if mac != "" && stb_ip_address != "" {
			status := auth.MacIpCheck(mac, stb_ip_address)
			token := services.GetAuthorizationToken(status, secret)
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{"authorization_token": token})
			return
		} else {
			level.Error(logger.Logger()).Log("Error", "MAC address and stb ip address should not be null")
			data := map[string]interface{}{
				"msg":        "Invalid Parameter",
				"error_code": 400,
			}
			w.WriteHeader(400)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)
			return
		}
	})

	r.Post("/create-nettv-migrate", func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		var post map[string]interface{}
		err := json.Unmarshal(reqBody, &post)
		if err != nil {
			fmt.Println(err)
		}
		auth.CreateNettvMigrate(post["username"].(string))
	})

	msg := fmt.Sprintf("Listening at http://localhost:%s", defaultPort)
	level.Info(logger.Logger()).Log("Status", msg)
	err = http.ListenAndServe(fmt.Sprintf(":%s", defaultPort), r)
	if err != http.ErrServerClosed {
		level.Error(logger.Logger()).Log("listen", err)
	}

}
