package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gookit/validate"
	"github.com/jsimonetti/pwscheme/md5crypt"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/authenticator/internal/repository"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/authenticator/internal/service"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/api"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type oldAuthController struct {
	authRepository       repository.AuthRepository
	nettvService         service.NettvService
	customerLoginService service.CustomerLogin
	apiWlink             api.ApiWlink
	logger               log.Logger
	oracle               *sql.DB
}

type OldAuthController interface {
	MacWiseUserCheck(c *gin.Context)
	AdlAuthentication(c *gin.Context)
	Authentication(c *gin.Context)
	OldSysTokenSync(c *gin.Context)
}

func NewOldSysAuth(r repository.AuthRepository, n service.NettvService, c service.CustomerLogin, a api.ApiWlink, l log.Logger, o *sql.DB) OldAuthController {
	return oldAuthController{
		authRepository:       r,
		nettvService:         n,
		customerLoginService: c,
		apiWlink:             a,
		logger:               l,
		oracle:               o,
	}
}

func (auth oldAuthController) MacWiseUserCheck(c *gin.Context) {
	var Payload types.MacwisePayload
	c.Bind(&Payload)
	payloadString, _ := json.Marshal(Payload)

	level.Info(auth.logger).Log("METHOD", "MacWiseUserCheck", "msg", "Mac-wise user check started", "payload", payloadString)
	v := validate.Struct(Payload)
	v.Validate()
	if !v.Errors.Empty() {
		res := types.OldSystemResponse{
			Status: 1,
			Data:   types.OldSystemMsg{Message: "Invalid Request"},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		level.Error(auth.logger).Log("METHOD", "MacWiseUserCheck", "msg", "Authentication Failed: invalid request", "payload", payloadString)
		return
	}

	splitMac := strings.Split(Payload.Mac, ":")
	joinMac := strings.Join(splitMac, "")
	macDetails, err := auth.authRepository.Search(bson.D{{"nettvsettopboxes", joinMac}})
	token, _ := auth.authRepository.GetToken(bson.D{{"stb_box_id", Payload.Mac}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fallbackData, status := auth.customerLoginService.RetailFallback(token.Username)
			if status {
				storeStatus, _ := auth.authRepository.StoreRetailCustomer(fallbackData)
				if storeStatus {
					userFallback, _ := auth.authRepository.Search(bson.D{{"username", token.Username}})
					fmt.Println("userfallback", userFallback)

					if !Contains(userFallback.STB_BOX_ID, joinMac) {
						res := types.OldSystemResponse{
							Status: 1,
							Data:   types.OldSystemMsg{Message: "Username and MAC maping not found"},
						}
						c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
						return
					}
					mappedStb := SetKeysForStbBoxes(userFallback.STB_BOX_ID)
					userFallback.MAPPED_STB_BOX_ID = mappedStb
					userFallback.Token = token.Token
					userFallback.Skipipcheck = strconv.FormatBool(token.SkipIpCheck.Valid)
					res := types.OldSystemResponse{
						Status: 0,
						Data:   types.OldSystemMsg{Message: userFallback},
					}
					c.JSON(http.StatusOK, []types.OldSystemResponse{res})
					return
				}
			}
			res := types.OldSystemResponse{
				Status: 1,
				Data:   types.OldSystemMsg{Message: "MAC Not Found"},
			}
			c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
			level.Error(auth.logger).Log("METHOD", "MacWiseUserCheck", "msg", "Mac-wise user check Failed: mac not found", "payload", payloadString)
			return
		}
	}

	if macDetails.Is_Nettv_Disabled == "Y" {
		res := types.OldSystemResponse{
			Status: 1,
			Data:   types.OldSystemMsg{Message: "Mac is disabled"},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		level.Error(auth.logger).Log("METHOD", "MacWiseUserCheck", "msg", "Mac-wise user check Failed: mac is disabled", "payload", payloadString)
		return
	}

	if macDetails.AccountStatus.Disable == "YES" {
		res := types.OldSystemResponse{
			Status: 1,
			Data:   types.OldSystemMsg{Message: "Sorry your account has been disabled."},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		level.Error(auth.logger).Log("METHOD", "MacWiseUserCheck", "msg", "Mac-wise user check Failed: account has beed disabled", "payload", payloadString)
		return
	}

	mappedStb := SetKeysForStbBoxes(macDetails.STB_BOX_ID)
	macDetails.MAPPED_STB_BOX_ID = mappedStb
	macDetails.AccountStatus.Password = ""
	macDetails.Token = token.Token
	macDetails.Skipipcheck = strconv.FormatBool(token.SkipIpCheck.Valid)
	res := types.OldSystemResponse{
		Status: 0,
		Data:   types.OldSystemMsg{Message: macDetails},
	}
	c.JSON(http.StatusOK, []types.OldSystemResponse{res})
	level.Error(auth.logger).Log("METHOD", "MacWiseUserCheck", "msg", "Mac-wise user check Failed: account has beed disabled", "payload", payloadString)
}

func (auth oldAuthController) AdlAuthentication(c *gin.Context) {
	var Payload types.AdlAuthenticationPayload
	c.Bind(&Payload)
	payloadString, _ := json.Marshal(Payload)
	level.Info(auth.logger).Log("METHOD", "AdlAuthentication", "msg", "AdlAuthentication started", "payload", payloadString)
	v := validate.Struct(Payload)
	v.Validate()
	if !v.Errors.Empty() {
		res := types.OldSystemResponse{
			Status: 1,
			Data:   types.OldSystemMsg{Message: "Invalid Request"},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		level.Error(auth.logger).Log("METHOD", "AdlAuthentication", "msg", "Authentication Failed: invalid request", "payload", payloadString)
		return
	}

	v.StringRule("username", "required")
	v.StringRule("password", "required")
	if !v.Validate() {
		res := types.OldSystemResponse{
			Status: 0,
			Data:   types.OldSystemMsg{Message: "Invalid Request"},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		level.Error(auth.logger).Log("METHOD", "AdlAuthentication", "msg", "AdlAuthentication Failed: invalid request", "payload", payloadString)
		return
	}

	response := auth.apiWlink.AdAuthenticate(Payload.Username, Payload.Password)
	if response.StatusCode == 200 {
		res := types.OldSystemResponse{
			Status: 0,
			Data:   types.OldSystemMsg{Message: response.Data},
		}
		c.JSON(http.StatusOK, []types.OldSystemResponse{res})
		level.Info(auth.logger).Log("METHOD", "AdlAuthentication", "msg", "AdlAuthentication success", "payload", payloadString)
		return
	}
}

func (auth oldAuthController) Authentication(c *gin.Context) {
	var Payload types.OldAuthPayload
	c.Bind(&Payload)
	payloadString, _ := json.Marshal(Payload)
	level.Info(auth.logger).Log("METHOD", "Authentication", "msg", "Authentication started", "payload", payloadString)

	v := validate.Struct(Payload)
	v.Validate()
	if !v.Errors.Empty() || (Payload.Password == "" && Payload.Token == "") {
		res := types.OldSystemResponse{
			Status: 1,
			Data:   types.OldSystemMsg{Message: "Invalid Request"},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		level.Error(auth.logger).Log("METHOD", "Authentication", "msg", "Authentication Failed: invalid request", "payload", payloadString)
		return
	}

	user, err := auth.authRepository.Search(bson.D{{"username", Payload.Username}})
	if err != nil || user.Is_Nettv_Disabled == "" {
		goTofallback(c, Payload, auth)
		return
	}
	splitMac := strings.Split(Payload.Mac, ":")
	joinMac := strings.Join(splitMac, "")
	if !Contains(user.STB_BOX_ID, joinMac) {
		res := types.OldSystemResponse{
			Status: 1,
			Data:   types.OldSystemMsg{Message: "Username and MAC maping not found"},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		return
	}

	if user.Is_Nettv_Disabled == "Y" {
		res := types.OldSystemResponse{
			Status: 1,
			Data:   types.OldSystemMsg{Message: "Account is disabled"},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		return
	}

	if user.AccountStatus.Disable == "YES" {
		res := types.OldSystemResponse{
			Status: 1,
			Data:   types.OldSystemMsg{Message: "Sorry your account has been disabled."},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		level.Error(auth.logger).Log("METHOD", "MacWiseUserCheck", "msg", "Mac-wise user check Failed: account has beed disabled", "payload", payloadString)
		return
	}

	token, _ := auth.authRepository.GetToken(bson.D{{"username", Payload.Username}, {"stb_box_id", Payload.Mac}})
	network, err := checkNetwork(Payload.Ip, user.AccountStatus.IpoeMac, token.SkipIpCheck.Valid, Payload.DeviceType)
	if err != nil {
		fmt.Println(err)
	}

	if !network {
		level.Error(auth.logger).Log("METHOD", "Authentication", "msg", err, "payload", payloadString)
		res := types.OldSystemResponse{
			Status: 1,
			Data:   types.OldSystemMsg{Message: "Network not allowed."},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		return
	}

	if Payload.Token != "" {
		if token.Token != Payload.Token {
			level.Error(auth.logger).Log("METHOD", "Authentication", "msg", "Authentication Error", "payload", payloadString)
			res := types.OldSystemResponse{
				Status: 1,
				Data:   types.OldSystemMsg{Message: "Authentication Error"},
			}
			c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
			return
		}
		level.Info(auth.logger).Log("METHOD", "Authentication", "msg", "Authenticate Successfully", "payload", payloadString)
		user.AccountStatus.Password = ""
		user.Token = token.Token
		user.Skipipcheck = strconv.FormatBool(token.SkipIpCheck.Valid)
		res := types.OldSystemResponse{
			Status: 0,
			Data:   types.OldSystemMsg{Message: user},
		}
		c.JSON(http.StatusOK, []types.OldSystemResponse{res})
		return
	}

	passwordHash := user.AccountStatus.Password
	verify, err := passwordVerify(Payload.Password, passwordHash)
	if err != nil {
		level.Error(auth.logger).Log("METHOD", "Authentication", "msg", err, "payload", payloadString)
		res := types.OldSystemResponse{
			Status: 1,
			Data:   types.OldSystemMsg{Message: "Authentication Error"},
		}
		c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
		return
	}

	if verify {
		token.Token = generateToken()
		auth.authRepository.SyncToken(token)
		level.Info(auth.logger).Log("METHOD", "Authentication", "msg", "Authenticate Successfully", "payload", payloadString)
		user.AccountStatus.Password = ""
		user.Token = token.Token
		user.Skipipcheck = strconv.FormatBool(token.SkipIpCheck.Valid)
		res := types.OldSystemResponse{
			Status: 0,
			Data:   types.OldSystemMsg{Message: user},
		}
		c.JSON(http.StatusOK, []types.OldSystemResponse{res})
		return
	}
}

func (auth oldAuthController) OldSysTokenSync(c *gin.Context) {
	t := time.Now()
	batchSize := 1000
	rows, err := auth.oracle.Query("SELECT username, token, stb_box_id, new_nettv_map FROM wlinkiptv_usermap WHERE username IN ( SELECT username FROM wlinkiptv_usermap GROUP BY username HAVING COUNT(*) = SUM( CASE WHEN new_nettv_map = 'N' THEN 1 ELSE 0 END ) AND COUNT(*) > 0 ) AND token IS NOT NULL")
	if err != nil {
		fmt.Println("err", err)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			fmt.Println("Can't close dataset: ", err)
		}
	}()
	for {
		data := make([]types.OldSysUser, 0, batchSize)
		for i := 0; i < batchSize && rows.Next(); i++ {
			var Username string
			var Token string
			var SkipIpCheck sql.NullString
			var StbBoxId string
			err = rows.Scan(&Username, &Token, &StbBoxId, &SkipIpCheck)
			if err != nil {
				fmt.Println("err", err)
			}
			userData := types.OldSysUser{Username: Username, SkipIpCheck: SkipIpCheck, Token: Token, StbBoxId: StbBoxId, Updated_At: time.Now()}
			data = append(data, userData)
		}
		if len(data) == 0 {
			break
		}

		documents := make([]interface{}, len(data))
		for i, person := range data {
			documents[i] = person
		}
		auth.authRepository.SyncTokenInBatch(documents)
	}
	level.Info(auth.logger).Log("METHOD", "OldSysTokenSync", "msg", "Successfully Synced", "time taken to sync", time.Since(t))
	c.JSON(http.StatusOK, gin.H{"message": "Successfully synced"})
}

func passwordVerify(password string, passwordHash string) (bool, error) {
	verify, err := md5crypt.Validate(password, `{MD5-CRYPT}`+passwordHash)
	if err != nil {
		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordHash))
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return verify, nil
}

func generateToken() string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 16)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func checkNetwork(ip string, ipoe string, skipIpCheck bool, deviceType string) (bool, error) {
	valid, err := regexp.MatchString("^10\\..*", ip)
	if err != nil {
		return false, err
	}
	if !valid && deviceType != "stb" && (ipoe == "" || !skipIpCheck) {
		return false, nil
	}
	return true, nil
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func goTofallback(c *gin.Context, Payload types.OldAuthPayload, auth oldAuthController) {
	fallbackData, status := auth.customerLoginService.RetailFallback(Payload.Username)
	if status {
		storeStatus, _ := auth.authRepository.StoreRetailCustomer(fallbackData)
		token, _ := auth.authRepository.GetToken(bson.D{{"username", Payload.Username}, {"stb_box_id", Payload.Mac}})

		if storeStatus {
			userFallback, _ := auth.authRepository.Search(bson.D{{"username", Payload.Username}})
			splitMac := strings.Split(Payload.Mac, ":")
			joinMac := strings.Join(splitMac, "")
			if !Contains(userFallback.STB_BOX_ID, joinMac) {
				res := types.OldSystemResponse{
					Status: 1,
					Data:   types.OldSystemMsg{Message: "Username and MAC maping not found"},
				}
				c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
				return
			}
			mappedStb := SetKeysForStbBoxes(userFallback.STB_BOX_ID)
			userFallback.MAPPED_STB_BOX_ID = mappedStb

			if userFallback.AccountStatus.Disable == "YES" {
				res := types.OldSystemResponse{
					Status: 1,
					Data:   types.OldSystemMsg{Message: "Sorry your account has been disabled."},
				}
				c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
				return
			}

			if userFallback.Is_Nettv_Disabled == "Y" {
				res := types.OldSystemResponse{
					Status: 1,
					Data:   types.OldSystemMsg{Message: "Account is disabled"},
				}
				c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
				return
			}

			if Payload.Token != "" {
				if token.Token != Payload.Token {
					level.Error(auth.logger).Log("METHOD", "goTofallback", "msg", "Fallback Incorrect Token", "payload", Payload.Username)
					res := types.OldSystemResponse{
						Status: 1,
						Data:   types.OldSystemMsg{Message: "Authentication Error"},
					}
					c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
					return
				}
				level.Info(auth.logger).Log("METHOD", "goTofallback", "msg", "Authenticate Successfully", "payload", Payload.Username)
				userFallback.AccountStatus.Password = ""
				userFallback.Token = token.Token
				userFallback.Skipipcheck = strconv.FormatBool(token.SkipIpCheck.Valid)
				res := types.OldSystemResponse{
					Status: 0,
					Data:   types.OldSystemMsg{Message: userFallback},
				}
				c.JSON(http.StatusOK, []types.OldSystemResponse{res})
				return
			}
			passwordHash := userFallback.AccountStatus.Password
			verify, err := passwordVerify(Payload.Password, passwordHash)
			if err != nil {
				level.Error(auth.logger).Log("METHOD", "goTofallback", "msg", err, "payload", Payload.Username)
				res := types.OldSystemResponse{
					Status: 1,
					Data:   types.OldSystemMsg{Message: "Authentication Error"},
				}
				c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
				return
			}

			if verify {
				token.Token = generateToken()
				auth.authRepository.SyncToken(token)
				level.Info(auth.logger).Log("METHOD", "goTofallback", "msg", "Authenticate Successfully", "payload", Payload.Username)
				userFallback.AccountStatus.Password = ""
				userFallback.Token = token.Token
				userFallback.Skipipcheck = strconv.FormatBool(token.SkipIpCheck.Valid)
				res := types.OldSystemResponse{
					Status: 0,
					Data:   types.OldSystemMsg{Message: userFallback},
				}
				c.JSON(http.StatusOK, []types.OldSystemResponse{res})
				return
			}
		}
	}

	level.Error(auth.logger).Log("METHOD", "goTofallback", "msg", "Fallback Username and MAC maping not found", "payload", Payload.Username)
	res := types.OldSystemResponse{
		Status: 1,
		Data:   types.OldSystemMsg{Message: "Username and MAC maping not found"},
	}
	c.AbortWithStatusJSON(http.StatusOK, []types.OldSystemResponse{res})
}
