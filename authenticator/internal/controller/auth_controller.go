package controller

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gookit/validate"
	"github.com/jsimonetti/pwscheme/md5crypt"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/authenticator/internal/repository"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/authenticator/internal/service"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type authController struct {
	authRepository       repository.AuthRepository
	nettvService         service.NettvService
	customerLoginService service.CustomerLogin
	logger               log.Logger
}

type AuthController interface {
	DefaultAuthenticate(c *gin.Context)
	CheckUserStatus(c *gin.Context)
	VerifyClient(c *gin.Context)
}

func New(r repository.AuthRepository, n service.NettvService, c service.CustomerLogin, l log.Logger) AuthController {
	return authController{
		authRepository:       r,
		nettvService:         n,
		customerLoginService: c,
		logger:               l,
	}
}

func (auth authController) DefaultAuthenticate(c *gin.Context) {
	var Payload map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&Payload)
	payloadString, _ := json.Marshal(Payload)
	v := validate.Map(Payload)
	v.StringRule("username", "required")
	v.StringRule("password", "required")
	level.Info(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Auhentication started", "payload", payloadString)
	if !v.Validate() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 422, "error_message": "invalid parameter"})
		level.Error(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Authentication Failed: invalid parameters", "payload", payloadString)
		return
	}
	checkCoperatePattern := CheckCorporatePattern(Payload["username"].(string))
	if checkCoperatePattern {
		level.Info(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Corporate User", "payload", payloadString)
		auth.handleCorporateAuth(Payload["username"].(string), Payload["password"].(string), c)
		return
	}
	val, err := auth.authRepository.GetUser(Payload["username"].(string))

	if err != nil {
		data, status := auth.customerLoginService.RetailFallback(Payload["username"].(string))
		if status {
			storeStatus, err := auth.authRepository.StoreRetailCustomer(data)
			if err != nil {
				fmt.Println("err", err)
			}
			if storeStatus {
				//fallback
				loginResp := auth.customerLoginService.GetCustomerLoginStatus(Payload["username"].(string), Payload["password"].(string))
				if loginResp.Success {
					if data.New_Nettv_Map != "Y" {
						//call authenticator
						level.Info(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Store in migrate table fallback", "payload", payloadString)
					}
					c.JSON(http.StatusOK, gin.H{"authenticaton": true, "data": data.AccountStatus})
					level.Info(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Authentication Successful fallback", "payload", payloadString, "Account-status", data.AccountStatus)
					return
				} else {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error_code": http.StatusUnprocessableEntity, "error_message": "Username or password didn't match"})
					level.Error(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Authentication Failed: Username or password didn't match fallback", "payload", payloadString)
					return
				}
			}
		} else {
			if err == mongo.ErrNoDocuments {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Username not found"})
				level.Error(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Username not found", "payload", payloadString)
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": err.Error()})
			level.Error(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", err.Error(), "payload", payloadString)
			return
		}
	}

	if val.AccountStatus.Username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Account status not present"})
		level.Info(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Authentication Failed: Account status not present", "payload", payloadString)
		return
	}

	if val.Is_Nettv_Disabled == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Nettv data not present"})
		level.Info(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Authentication Failed: Nettv data not present", "payload", payloadString)
		return
	}

	if val.Is_Nettv_Disabled == "Y" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Nettv Mapping Disabled"})
		level.Info(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Authentication Failed: Nettv Mapping Disabled", "payload", payloadString)
		return
	}
	var as types.AccountStatus
	if val.AccountStatus != as {
		checked, err := verifyCrypt(Payload["password"].(string), val.AccountStatus.Password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error_code": http.StatusUnprocessableEntity, "error_message": "Username or password didn't match"})
			level.Error(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Authentication Failed: Username or password didn't match", "payload", payloadString)
			return
		}

		if checked {
			if val.New_Nettv_Map != "Y" {
				//call authenticator
				level.Info(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Store in migrate table", "payload", payloadString)
			}
			val.AccountStatus.Password = ""
			c.JSON(http.StatusOK, gin.H{"authenticaton": true, "data": val.AccountStatus})
			level.Info(auth.logger).Log("METHOD", "DefaultAuthenticate", "msg", "Authentication Successful", "payload", payloadString, "Account-status", val.AccountStatus)
			return
		}
	}
}

func (auth authController) VerifyClient(c *gin.Context) {
	var Payload map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&Payload)
	payloadString, _ := json.Marshal(Payload)
	level.Info(auth.logger).Log("METHOD", "VerifyClient", "msg", "Verify client initiated", "payload", payloadString)
	v := validate.Map(Payload)
	v.StringRule("username", "required")
	if !v.Validate() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 422, "error_message": "invalid parameter"})
		level.Error(auth.logger).Log("METHOD", "VerifyClient", "msg", "Authentication Failed: invalid parameters", "payload", payloadString)
		return
	}
	username := Payload["username"].(string)
	checkCoperatePattern := CheckCorporatePattern(username)
	if checkCoperatePattern {
		type corporateReturnVal struct {
			Status int16               `json:"status"`
			Data   types.CorporateData `json:"data"`
		}
		val, err := auth.authRepository.GetCorporateUser(username)
		if err != nil {
			fallbackData := auth.customerLoginService.CorporateFallback(username)
			if fallbackData.NewNettvMap == "" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Username not found"})
				level.Error(auth.logger).Log("METHOD", "VerifyClient", "msg", "Username not found", "username", username)
				return
			}
			_, err := auth.authRepository.StoreCorporateCustomerFallback(fallbackData)
			if err != nil {
				level.Error(auth.logger).Log("METHOD", "VerifyClient", "msg", "Failed while storing corporate fallback", "username", username)
			}
			corporateFallbackReturnData := corporateReturnVal{Status: 200, Data: fallbackData}
			c.JSON(http.StatusOK, []corporateReturnVal{corporateFallbackReturnData})
			return
		}
		corporateFallbackReturnData := corporateReturnVal{Status: 200, Data: val}
		c.JSON(http.StatusOK, []corporateReturnVal{corporateFallbackReturnData})
		return
	}

	val, err := auth.authRepository.GetUser(Payload["username"].(string))
	type returnVal struct {
		Status int16      `json:"status"`
		Data   types.Data `json:"data"`
	}
	if err != nil {
		data, status := auth.customerLoginService.RetailFallback(Payload["username"].(string))
		if status {
			storeStatus, err := auth.authRepository.StoreRetailCustomer(data)
			if err != nil {
				fmt.Println("err", err)
			}
			if storeStatus {
				mappedStb := SetKeysForStbBoxes(data.STB_BOX_ID)
				data.MAPPED_STB_BOX_ID = mappedStb
				data.AccountStatus.Password = ""
				fallbackReturnData := returnVal{Status: 200, Data: data}
				c.JSON(http.StatusOK, []returnVal{fallbackReturnData})
				level.Info(auth.logger).Log("METHOD", "VerifyClient", "msg", "Authentication Successful fallback", "payload", payloadString, "Data", data)
				return
			}
		} else {
			if err == mongo.ErrNoDocuments {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Username not found"})
				level.Error(auth.logger).Log("METHOD", "VerifyClient", "msg", "Username not found", "payload", payloadString)
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": err.Error()})
			level.Error(auth.logger).Log("METHOD", "VerifyClient", "msg", err.Error(), "payload", payloadString)
			return
		}
	}
	if val.Is_Nettv_Disabled == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Nettv data not present"})
		level.Error(auth.logger).Log("METHOD", "VerifyClient", "msg", "Authentication Failed: Nettv data not present", "payload", payloadString)
		return
	}

	if val.Is_Nettv_Disabled == "Y" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Nettv Mapping Disabled"})
		level.Error(auth.logger).Log("METHOD", "VerifyClient", "msg", "Authentication Failed: Nettv Mapping Disabled", "payload", payloadString)
		return
	}
	mappedStb := SetKeysForStbBoxes(val.STB_BOX_ID)
	val.MAPPED_STB_BOX_ID = mappedStb
	val.AccountStatus.Password = ""
	returnData := returnVal{Status: 200, Data: val}
	c.JSON(http.StatusOK, []returnVal{returnData})
	jsonval, _ := json.Marshal(val)
	level.Info(auth.logger).Log("METHOD", "VerifyClient", "msg", "Authentication Successful", "payload", payloadString, "Account-status", string(jsonval))
}

func (auth authController) CheckUserStatus(c *gin.Context) {
	var Payload map[string]interface{}
	json.NewDecoder(c.Request.Body).Decode(&Payload)
	payloadString, _ := json.Marshal(Payload)
	v := validate.Map(Payload)
	v.StringRule("username", "required")
	v.StringRule("password", "required")
	if !v.Validate() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 422, "error_message": "invalid parameter"})
		level.Error(auth.logger).Log("METHOD", "CheckUserStatus", "msg", "Authentication Failed: invalid parameter", "payload", payloadString)
		return
	}
	username := Payload["username"].(string)
	password := Payload["password"].(string)
	checkCoperatePattern := CheckCorporatePattern(username)
	if checkCoperatePattern {
		level.Info(auth.logger).Log("METHOD", "CheckUserStatus", "msg", "Corporate: check user status", "payload", payloadString)
		md5Password := fmt.Sprintf("%x", md5.Sum([]byte(password)))
		splitUsername := strings.Split(username, "-")
		customerCode := splitUsername[0]
		checkCA := auth.corporateAuthentcation(customerCode, md5Password, c)
		if checkCA {
			c.JSON(http.StatusOK, gin.H{"authentication": true, "code": 200})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": http.StatusUnprocessableEntity, "error_message": "Invalid username or password"})
			level.Error(auth.logger).Log("METHOD", "CheckUserStatus", "msg", "Corporate: Invalid username or password", "payload", payloadString)
		}
		return
	}
	val, err := auth.authRepository.GetUser(username)
	if err != nil {
		data, status := auth.customerLoginService.RetailFallback(Payload["username"].(string))
		if status {
			storeStatus, err := auth.authRepository.StoreRetailCustomer(data)
			if err != nil {
				level.Error(auth.logger).Log("METHOD", "CheckUserStatus", "msg", "Authentication Failed: Fallback failed during store", "payload", payloadString)
			}
			if storeStatus {
				//fallback
				loginResp := auth.customerLoginService.GetCustomerLoginStatus(Payload["username"].(string), Payload["password"].(string))
				if loginResp.Success {
					c.JSON(http.StatusOK, gin.H{"authenticaton": true, "code": 200})
					level.Info(auth.logger).Log("METHOD", "CheckUserStatus", "msg", "Authentication Successful fallback", "payload", payloadString)
					return
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error_code": http.StatusUnprocessableEntity, "error_message": "Username or password didn't match"})
				level.Error(auth.logger).Log("METHOD", "CheckUserStatus", "msg", "Authentication Failed: Username or password didn't match fallback", "payload", payloadString)
				return
			}
		} else {
			if err == mongo.ErrNoDocuments {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Username not found"})
				level.Error(auth.logger).Log("METHOD", "CheckUserStatus", "msg", "Username not found", "payload", payloadString)
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": err.Error()})
			level.Error(auth.logger).Log("METHOD", "CheckUserStatus", "msg", err.Error(), "payload", payloadString)
			return
		}
	}

	checked, err := verifyCrypt(Payload["password"].(string), val.AccountStatus.Password)
	if err != nil {
		fmt.Println(err)
	}
	if checked {
		c.JSON(http.StatusOK, gin.H{"authenticaton": true, "code": 200})
		level.Info(auth.logger).Log("METHOD", "CheckUserStatus", "msg", "Authentication Successful", "payload", payloadString)
		return
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error_code": http.StatusUnprocessableEntity, "error_message": "Username or password didn't match"})
	level.Error(auth.logger).Log("METHOD", "CheckUserStatus", "msg", "Authentication Failed: Username or password didn't match", "payload", payloadString)
}

func CheckCorporatePattern(reg string) bool {
	matched, _ := regexp.MatchString("^([A-Z0-9]{4,5})-([A-Z]{1}[0-9]{2})-([0-9]{2})", reg)
	return matched
}

func SetKeysForStbBoxes(stbBoxes []string) map[string]interface{} {
	stbs := make(map[string]interface{})
	for i, v := range stbBoxes {
		key := fmt.Sprintf("id%+v", i)
		stbs[key] = v
	}
	return stbs
}

func (auth authController) handleCorporateAuth(username string, password string, c *gin.Context) {
	level.Info(auth.logger).Log("METHOD", "handleCorporateAuth", "msg", "handle corporate auth start", "username", username)
	val, err := auth.authRepository.GetCorporateUser(username)
	if err != nil {
		fallbackData := auth.customerLoginService.CorporateFallback(username)
		if fallbackData.NewNettvMap == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Username not found"})
			level.Error(auth.logger).Log("METHOD", "handleCorporateAuth", "msg", "Username not found", "username", username)
			return
		}
		_, err := auth.authRepository.StoreCorporateCustomerFallback(fallbackData)
		if err != nil {
			level.Error(auth.logger).Log("METHOD", "handleCorporateAuth", "msg", "Failed while storing corporate fallback", "username", username)
		}
		checkCA := auth.corporateAuthentcation(fallbackData.CustomerCode, fallbackData.Password, c)
		if checkCA {
			returnData := types.CorporateReturnData{}
			returnData.Username = fallbackData.Username
			returnData.Address = fallbackData.CustomerInfo.Address
			returnData.Customer = fallbackData.CustomerInfo.Customer
			returnData.DISABLE = fallbackData.AccountStatus
			returnData.Email = fallbackData.CustomerInfo.Email
			returnData.Phone1 = fallbackData.CustomerInfo.Phone1
			returnData.NEW_NETTV_MAP = fallbackData.NewNettvMap
			c.JSON(http.StatusOK, gin.H{"authentication": true, "data": returnData})
			level.Info(auth.logger).Log("METHOD", "handleCorporateAuth", "msg", "Corporate authentication success fallback", "username", username)
			return
		} else {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error_code": 422, "error_message": "Invalid Authentication."})
			level.Error(auth.logger).Log("METHOD", "handleCorporateAuth", "msg", "Invalid Authentication fallback.", "username", username)
			return
		}
	}

	if val.IsNettvDisabled == "Y" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error_code": 400, "error_message": "Nettv Mapping Disabled"})
		level.Error(auth.logger).Log("METHOD", "handleCorporateAuth", "msg", "Nettv Mapping Disabled", "username", username)
		return
	}

	md5Password := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	splitUsername := strings.Split(username, "-")
	customerCode := splitUsername[0]
	checkCA := auth.corporateAuthentcation(customerCode, md5Password, c)
	if checkCA {
		returnData := types.CorporateReturnData{}
		returnData.Username = val.Username
		returnData.Address = val.CustomerInfo.Address
		returnData.Customer = val.CustomerInfo.Customer
		returnData.DISABLE = val.AccountStatus
		returnData.Email = val.CustomerInfo.Email
		returnData.Phone1 = val.CustomerInfo.Phone1
		returnData.NEW_NETTV_MAP = val.NewNettvMap
		c.JSON(http.StatusOK, gin.H{"authentication": true, "data": returnData})
		level.Info(auth.logger).Log("METHOD", "handleCorporateAuth", "msg", "Corporate authentication success", "username", username)
	} else {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error_code": 422, "error_message": "Invalid Authentication."})
		level.Error(auth.logger).Log("METHOD", "handleCorporateAuth", "msg", "Invalid Authentication.", "username", username)
	}
}

func (auth authController) corporateAuthentcation(customerCode string, password string, c *gin.Context) bool {
	level.Info(auth.logger).Log("METHOD", "corporateAuthentcation", "msg", "Authentication started", "username", customerCode)
	codeData, err := auth.authRepository.GetCustomerCode(customerCode)
	if err != nil {
		data, statusCode := auth.customerLoginService.CorporateCustomerCodeFallback(customerCode)
		if statusCode {
			var customerCodePayload = types.CorporateCustomerCodeData{}
			customerCodePayload.CustomerCode = data.Code
			customerCodePayload.AccountStatus = "Y"
			customerCodePayload.Password = data.Password
			customerCodePayload.Updated_At = time.Now()
			status, _ := auth.authRepository.StoreCorporateCustomerCode(customerCodePayload)
			if status {
				if data.Password == password {
					return true
				}
			}
		} else {
			level.Info(auth.logger).Log("METHOD", "corporateAuthentcation", "msg", err.Error(), "username", customerCode)
			return false
		}
	}
	if codeData.Password == password {
		return true
	}
	return false
}

func verifyCrypt(password string, hashedPassword string) (bool, error) {
	verify, err := md5crypt.Validate(password, `{MD5-CRYPT}`+hashedPassword)
	if err != nil {
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return verify, nil
}
