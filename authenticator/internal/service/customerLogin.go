package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/api"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
)

type CustomerLogin interface {
	GetCustomerLoginStatus(username string, password string) CustomerLoginResponse
	RetailFallback(username string) (types.Data, bool)
	CorporateFallback(username string) types.CorporateData
	CorporateCustomerCodeFallback(username string) (types.CustomerCodeFallbackResponseData, bool)
}

type CustomerLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type customerlogin struct {
	logger log.Logger
}

func NewCustomerLogin(logger log.Logger) CustomerLogin {
	return &customerlogin{
		logger: logger,
	}
}

const CustomerLoginUrl = "https://customer-login-api.wlink.com.np/"

func (n *customerlogin) GetCustomerLoginStatus(username string, password string) CustomerLoginResponse {
	level.Info(n.logger).Log("METHOD", "GetCustomerLoginStatus", "msg", "Cutomer Login started ")
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)
	req, err := http.NewRequest("POST", CustomerLoginUrl, strings.NewReader(form.Encode()))
	if err != nil {
		level.Error(n.logger).Log("METHOD", "GetCustomerLoginStatus", "msg", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth("customer", "b4f8WGkep4RGW8EH")
	resp, err := client.Do(req)
	if err != nil {
		level.Error(n.logger).Log("METHOD", "GetCustomerLoginStatus", "msg", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		level.Error(n.logger).Log(err)
	}
	defer resp.Body.Close()
	level.Info(n.logger).Log("METHOD", "GetCustomerLoginStatus", "statusCode", resp.StatusCode, "response", string(body))
	var clr CustomerLoginResponse
	json.Unmarshal(body, &clr)
	level.Info(n.logger).Log("METHOD", "GetCustomerLoginStatus", "msg", "Cutomer Login success", "RESPONSE", string(body))
	return clr
}

func (n *customerlogin) RetailFallback(username string) (types.Data, bool) {
	customerService := api.WlinkService(n.logger)
	nettvService := api.NewNettvService(n.logger)
	userStatus := customerService.GetUserStatus(username)
	if userStatus["user_name"] == nil {
		return types.Data{}, false
	}
	userBandwidth := customerService.GetCustomerAccountDetail(userStatus["plan_category_id"].(string))
	upBandwidth := int64(userBandwidth["up_bw"].(float64)) / 1024
	downBandwidth := int64(userBandwidth["down_bw"].(float64)) / 1024
	minUpBandwidth := userBandwidth["min_up_bw"]
	bandwidth := fmt.Sprintf("%d/%d Mbps", upBandwidth, downBandwidth)
	fallbackspeed := "NO"
	graceStatus := userStatus["grace_status"].(string)
	vol := "EXP"
	if minUpBandwidth != nil {
		vol = "VOL"
		minDownBandwidth := userBandwidth["min_down_bw"]
		fallbackspeed = fmt.Sprintf("%s/%s Kbps (Running on Fallback speed)", minUpBandwidth, minDownBandwidth)
	}
	disable := userStatus["disable"].(string)
	acStatus := "Active"
	if disable == "N" {
		disable = "NO"
		if vol == "VOL" {
			volrem := userStatus["volume_remaining"].(string)
			vr, err := strconv.Atoi(volrem)
			if err != nil {
				fmt.Println(err)
			}
			roundvol := math.Round((((float64(vr) / 1024) / 1024) / 1024))
			if roundvol <= 0 {
				acStatus = "Fallback"
			}
		} else if vol == "EXP" {
			if graceStatus == "Y" {
				acStatus = "Grace Period"
				graceStatus = "YES"
			} else {
				acStatus = "Active"
				graceStatus = "NO"
			}
		}
	} else if disable == "Y" {
		disable = "YES"
		acStatus = "Disabled"
	}

	accountStatus := types.AccountStatus{
		DaysRemaining:       strings.TrimSpace(userStatus["days_remaining"].(string)),
		GraceStatus:         strings.TrimSpace(graceStatus),
		Disable:             strings.TrimSpace(disable),
		CurrentBandwidth:    strings.TrimSpace(bandwidth),
		ExpiryDate:          strings.TrimSpace(userStatus["expiry_date"].(string)),
		MinuteLeft:          strings.TrimSpace(userStatus["minutes_remaining"].(string)),
		SubscribedBandwidth: strings.TrimSpace(bandwidth),
		VolumeLeft:          strings.TrimSpace(userStatus["volume_remaining"].(string)),
		PayPlan:             strings.TrimSpace(userStatus["pay_plan_id"].(string)),
		Username:            strings.TrimSpace(username),
		FallbackSpeed:       strings.TrimSpace(fallbackspeed),
		AccountStatus:       strings.TrimSpace(acStatus),
		Balance:             "0",
		PlanName:            userStatus["pay_plan_name"].(string) + " " + userStatus["plan_category_name"].(string),
		// LastOnline:    Cspayload.Online_Date,
	}
	customerInfo := types.CustomerInfo{}
	customer := customerService.GetCustomerInfo(username)
	customerInfo.Username = username
	customerInfo.Address = customer["address"].(string)
	if customer["primary_email_address"] != nil {
		customerInfo.PrimaryEmailAddress = strings.TrimSpace(customer["primary_email_address"].(string))
	}
	if customer["primary_mobile_number"] != nil {
		customerInfo.PrimaryMobileNumber = strings.TrimSpace(customer["primary_mobile_number"].(string))
	}
	customerInfo.FullName = customer["name"].(string)
	if customer["primary_phone_line"] != nil {
		customerInfo.PHONE1 = customer["primary_phone_line"].(string)
	}
	allData := types.Data{}
	customerPassword := customerService.GetCustomerDetailsWithPassword(username)
	if customerPassword["pass_word_admin"] != nil {
		accountStatus.Password = customerPassword["pass_word_admin"].(string)
	}
	if customerPassword["ipoe_mac"] != nil {
		accountStatus.IpoeMac = customerPassword["ipoe_mac"].(string)
	}
	allData.Username = accountStatus.Username
	allData.AccountStatus = accountStatus
	allData.CustomerInfo = customerInfo
	allData.Data_Type = "retail-nettv"
	allData.Updated_At = time.Now()
	nettv := nettvService.GetCustomerNettvDetail(username)
	if nettv != nil {
		stbbox, newSystem := nettvService.FilterNettvData(nettv)
		if stbbox != nil {
			allData.New_Nettv_Map = newSystem
			allData.Is_Nettv_Disabled = "N"
			allData.STB_BOX_ID = stbbox
			allData.Session_Count = "0"
			allData.Skipipcheck = "false"
		}
	}
	return allData, true
}

func (n *customerlogin) CorporateFallback(username string) types.CorporateData {
	var CorporateCustomerInfo = types.CorporateCustomerInfo{}
	corporateService := api.CorporateService(n.logger)
	csData := corporateService.GetCorporateUserDetail(username)
	var customerDetailFallback types.CorporateCustFallbackResponse
	val, _ := json.Marshal(csData)
	json.Unmarshal(val, &customerDetailFallback)
	cd := types.CorporateData{}
	if len(customerDetailFallback.Data) > 0 {
		c := customerDetailFallback.Data[0]
		CorporateCustomerInfo.Address = c.Address
		CorporateCustomerInfo.Customer = c.Customer
		CorporateCustomerInfo.Email = c.Email
		CorporateCustomerInfo.Mobile = c.Mobile
		CorporateCustomerInfo.Phone1 = c.Phone1
		CorporateCustomerInfo.Username = c.Username
		cd.Username = c.Username
		cd.CustomerCode = c.CustomerCode
		cd.CustomerInfo = CorporateCustomerInfo
		cd.AccountStatus = c.AccountStatus
		cd.Updated_At = time.Now()
		cd.Data_type = "corporate-nettv"
		cd.Password = c.Password
		cd.SessionCount = "0"
		cpNettvService := api.CorporateNettvService(n.logger)
		stbs, sys, exists := cpNettvService.GetCorporateStbBoxes(username)
		cd.StbBoxId = stbs
		cd.NewNettvMap = sys
		cd.IsNettvDisabled = "N"
		if !exists {
			cd.IsNettvDisabled = "Y"
		}
		return cd
	} else {
		return cd
	}
}

func (n *customerlogin) CorporateCustomerCodeFallback(username string) (types.CustomerCodeFallbackResponseData, bool) {
	corporateService := api.CorporateService(n.logger)
	customerCodeData := corporateService.GetCustomerCodeDetail(username)
	var codeFallbackData types.CustomerCodeFallbackResponse
	val, _ := json.Marshal(customerCodeData)
	json.Unmarshal(val, &codeFallbackData)
	level.Info(n.logger).Log("METHOD", "CorporateCustomerCodeFallback", "msg", "fallback successful", "username", username)
	if len(codeFallbackData.Data) > 0 {
		return codeFallbackData.Data[0], true
	} else {
		return types.CustomerCodeFallbackResponseData{}, false
	}
}
