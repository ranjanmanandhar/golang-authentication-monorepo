package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/rabbitmq/amqp091-go"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/api"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	logger                log.Logger
	mongo                 MongoDBService
	redis                 RedisService
	swlink                ServiceWlink
	nettvService          NettvService
	corporateService      api.ServiceCorporate
	corporateNettvService api.CorporateNettv
}

func NewService(logger log.Logger, mongoService MongoDBService, redisService RedisService, serviceWlink ServiceWlink, nettvService NettvService, corporateService api.ServiceCorporate, corporateNettvService api.CorporateNettv) service {
	return service{
		logger:                logger,
		mongo:                 mongoService,
		redis:                 redisService,
		swlink:                serviceWlink,
		nettvService:          nettvService,
		corporateService:      corporateService,
		corporateNettvService: corporateNettvService,
	}
}

func (s service) ProcessQueue(ctx context.Context, d amqp091.Delivery) {
	var body map[string]interface{}
	err := json.Unmarshal(d.Body, &body)
	if err != nil {
		level.Error(s.logger).Log("Error", err)
	}
	jsonBody, _ := json.Marshal(body)
	level.Info(s.logger).Log("METHOD", "ProcessQueue", "queue messsage", jsonBody)
	if body["type"] == "corporate" {
		s.HandleCorporateCustomer(ctx, body)
	} else if body["type"] == "corporate-nettv" {
		s.HandleCorporateNettv(ctx, body)
	} else if body["type"] == "nettv" {
		s.HandleRetailNettv(ctx, body)
	} else if body["schema"] == "customer" {
		s.HandleRetailCustomer(ctx, body)
	} else if body["schema"] == "customerinfo" {
		s.HandleRetailCustomerInfo(ctx, body)
	}
}

func (s service) HandleCorporateNettv(ctx context.Context, body map[string]interface{}) {
	var filter bson.M
	var CorporateData types.CorporateNettvData
	CorporatePayload := types.CorporateNettvPayload{}
	val, _ := json.Marshal(body["data"])
	json.Unmarshal(val, &CorporatePayload)
	filter = bson.M{
		"username": CorporatePayload.Username,
	}
	CorporateData.Username = CorporatePayload.Username
	CorporateData.StbBoxId = CorporatePayload.StbBoxIds
	prevData, err := s.mongo.FindOne(ctx, filter)
	if prevData["stbboxid"] == nil && (body["action"] == "detach-cnd" || body["action"] == "detach-ebill") {
		level.Error(s.logger).Log("msg", "No stbs to detach")
		return
	}
	if err == nil && prevData["stbboxid"] != nil {
		stbbox := GetCorporateSettopboxes(prevData["stbboxid"].(primitive.A), CorporateData.StbBoxId, body["action"].(string))
		CorporateData.StbBoxId = stbbox
	}
	s.mongo.UpdateOrInsertCorporateNettv(ctx, filter, CorporateData)
	s.redis.DelFromRedis(ctx, CorporateData.Username)
}

func (s service) HandleCorporateCustomer(ctx context.Context, body map[string]interface{}) {
	if body["action"] == "create_customer_code" || body["action"] == "update_password" {
		CorporateData := types.CorporateCustomerCodeData{}
		val, _ := json.Marshal(body["data"])
		json.Unmarshal(val, &CorporateData)
		var customerCodeFilter bson.M
		customerCodeFilter = bson.M{
			"customer_code": CorporateData.CustomerCode,
		}
		CorporateData.AccountStatus = "N"
		s.mongo.UpdateOrInsertCorporateCodes(ctx, customerCodeFilter, CorporateData)
	}

	if body["action"] == "delete_customer_code" {
		CorporateData := types.CorporateCustomerCodeData{}
		val, _ := json.Marshal(body["data"])
		json.Unmarshal(val, &CorporateData)
		var customerCodeFilter bson.M
		customerCodeFilter = bson.M{
			"customer_code": CorporateData.CustomerCode,
		}
		CorporateData.AccountStatus = "Y"
		s.mongo.UpdateOrInsertCorporateCodes(ctx, customerCodeFilter, CorporateData)
	}

	if body["action"] == "create_circuit_id" || body["action"] == "update_cid" {
		corporatePayload := types.CorporateNettvPayload{}
		val, _ := json.Marshal(body["data"])
		json.Unmarshal(val, &corporatePayload)

		CorporateData := types.CustomerInfo{}
		val, _ = json.Marshal(body["data"])
		json.Unmarshal(val, &CorporateData)
		var corporateFilter bson.M
		corporateFilter = bson.M{
			"username": CorporateData.Username,
		}
		var CustomerInfo types.CorporateCustomerInfo
		json.Unmarshal(val, &CustomerInfo)
		var CircuitData types.CircuitData
		CircuitData.Username = CorporateData.Username
		CircuitData.CustomerInfo = CustomerInfo
		CircuitData.AccountStatus = corporatePayload.AccountStatus
		CircuitData.CustomerCode = corporatePayload.CustomerCode
		CircuitData.Updated_At = time.Now()
		CircuitData.Data_type = "corporate-nettv"
		if body["action"] == "create_circuit_id" {
			CircuitData.CircuitIdStatus = "Active"
			CircuitData.AccountStatus = "N"
		}
		s.mongo.UpdateOrInsertCorporate(ctx, corporateFilter, CircuitData)
		s.redis.DelFromRedis(ctx, corporatePayload.Username)
	}
}

func (s service) HandleRetailNettv(ctx context.Context, body map[string]interface{}) {
	var filter bson.M
	var NettvData types.RetailNettvData
	NettvPayload := types.NettvPayload{}
	val, _ := json.Marshal(body)
	json.Unmarshal(val, &NettvPayload)
	filter = bson.M{
		"username": NettvPayload.Username,
	}
	NettvData.Username = NettvPayload.Username
	NettvData.New_Nettv_Map = NettvPayload.NewNetvMap
	NettvData.STB_BOX_ID = []string{NettvPayload.StbBoxIds}
	prevData, err := s.mongo.FindOne(ctx, filter)
	if prevData["nettvsettopboxes"] == nil && NettvPayload.DataType == "detach" {
		level.Error(s.logger).Log("msg", "No stbs to detach")
		return
	}
	customerInfo := types.CustomerInfo{}
	if prevData["customerInfo"] == nil {
		customer := s.swlink.GetCustomerInfo(NettvPayload.Username)
		customerInfo.Username = NettvPayload.Username
		customerInfo.Address = customer["address"].(string)
		if customer["primary_email_address"] != nil {
			customerInfo.PrimaryEmailAddress = customer["primary_email_address"].(string)
		}
		if customer["primary_mobile_number"] != nil {
			customerInfo.PrimaryMobileNumber = customer["primary_mobile_number"].(string)
		}
		customerInfo.FullName = customer["name"].(string)
		if customer["primary_phone_line"] != nil {
			customerInfo.PHONE1 = customer["primary_phone_line"].(string)
		}
	}
	if err == nil && prevData["nettvsettopboxes"] != nil {
		NettvData.STB_BOX_ID = Getprevsettopboxes(prevData["nettvsettopboxes"].(primitive.A), NettvPayload.StbBoxIds, NettvPayload.DataType)
		level.Info(s.logger).Log("msg", "successfully fetched perviously stored stb boxes", "previous setbboxes", NettvData.STB_BOX_ID)
	}
	if err != nil {
		level.Error(s.logger).Log("msg", err)
	}
	NettvData.Is_Nettv_Disabled = "N"
	NettvData.New_Nettv_Map = "Y"
	NettvData.Skip_Ip_Check = "false"
	NettvData.Session_Count = "0"
	NettvData.Updated_At = time.Now()
	NettvData.Data_type = "retail-nettv"
	s.mongo.UpdateOrInsertNettv(ctx, filter, NettvData)
	s.redis.DelFromRedis(ctx, NettvPayload.Username)
}

func (s service) HandleRetailCustomerInfo(ctx context.Context, body map[string]interface{}) {
	level.Info(s.logger).Log("method", "HandleRetailCustomerInfo")
	var filter bson.M
	Cspayload := types.NewCspayload{}
	val, _ := json.Marshal(body["data"])
	json.Unmarshal(val, &Cspayload)
	filter = bson.M{
		"username": Cspayload.Machine_Name,
	}
	customerInfo := types.CustomerInfo{
		Username:            Cspayload.Username,
		Address:             Cspayload.Per_Address,
		PrimaryEmailAddress: Cspayload.Per_Cont_Email_Primary,
		PrimaryMobileNumber: Cspayload.Per_Cont_Mobile,
		FullName:            Cspayload.Per_Client_Name,
		PHONE1:              Cspayload.Per_Cont_Mobile,
		PHONE2:              Cspayload.Per_Cont_Mobile,
	}
	level.Info(s.logger).Log("msg", "successfully mapped customer info")
	custdata := types.CustomerInfoData{
		Username:     Cspayload.Username,
		CustomerInfo: customerInfo,
		Updated_At:   time.Now(),
		Data_type:    "retail-nettv",
	}

	s.mongo.UpdateOrInsertCustomerInfo(ctx, filter, custdata)
	level.Info(s.logger).Log("msg", "successfully updated data from customer service", "payload", custdata)

	s.redis.DelFromRedis(ctx, Cspayload.Username)
}

func (s service) HandleRetailCustomer(ctx context.Context, body map[string]interface{}) {
	level.Info(s.logger).Log("method", "HandleRetailCustomer")

	var filter bson.M
	Cspayload := types.NewCspayload{}
	val, _ := json.Marshal(body["data"])
	json.Unmarshal(val, &Cspayload)
	filter = bson.M{
		"username": Cspayload.Username,
	}
	// prevData, err := s.mongo.FindOne(ctx, filter)
	// if err != nil {
	// 	level.Error(s.logger).Log("msg", err)
	// }
	userStatus := s.swlink.GetUserStatus(Cspayload.Username)
	userBandwidth := s.swlink.GetCustomerAccountDetail(userStatus["plan_category_id"].(string))
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
		Username:            strings.TrimSpace(Cspayload.Username),
		FallbackSpeed:       strings.TrimSpace(fallbackspeed),
		AccountStatus:       strings.TrimSpace(acStatus),
		Balance:             "0",
		PlanName:            userStatus["pay_plan_name"].(string) + " " + userStatus["plan_category_name"].(string),
		// LastOnline:    Cspayload.Online_Date,
	}

	level.Info(s.logger).Log("msg", "successfully mapped customer account status")

	customerInfo := types.CustomerInfo{}
	customer := s.swlink.GetCustomerInfo(Cspayload.Username)
	customerInfo.Username = Cspayload.Username
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
	customerPassword := s.swlink.GetCustomerDetailsWithPassword(Cspayload.Username)
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
	nettv := s.nettvService.GetCustomerNettvDetail(Cspayload.Username)
	if nettv != nil {
		stbbox, newSystem := FilterNettvData(nettv)
		if stbbox != nil {
			allData.New_Nettv_Map = newSystem
			allData.Is_Nettv_Disabled = "N"
			allData.STB_BOX_ID = stbbox
			allData.Session_Count = "0"
			allData.Skipipcheck = "false"
		}
	}
	s.mongo.UpdateOrInsertCustomer(ctx, filter, allData)
	level.Info(s.logger).Log("msg", "successfully updated data from customer service", "payload", allData)
	s.redis.DelFromRedis(ctx, Cspayload.Username)
}

func FilterNettvData(nettv []map[string]interface{}) ([]string, string) {
	var stbbox []string
	newSystem := "N"
	for _, val := range nettv {
		stbbox = append(stbbox, val["stb_box_id"].(string))
		if val["new_system"] == true {
			newSystem = "Y"
		}
	}
	return stbbox, newSystem
}

func Getprevsettopboxes(setopBoxes primitive.A, stbId string, action string) []string {
	check := Contains(setopBoxes, stbId)
	var allStbs []string
	if !check {
		allStbs = append(allStbs, stbId)
	}
	for _, val := range setopBoxes {
		if action == "detach" && val.(string) == stbId {
			continue
		}
		allStbs = append(allStbs, val.(string))
	}
	return allStbs
}

func GetCorporateSettopboxes(setopBoxes primitive.A, stbIds []string, action string) []string {
	var filtered []string
	for _, stbId := range stbIds {
		check := Contains(setopBoxes, stbId)
		if !check {
			filtered = append(filtered, stbId)
		}
	}
	for _, val := range setopBoxes {
		if (action == "detach-cdn" || action == "detach-ebill") && Contains(setopBoxes, val.(string)) {
			continue
		}
		filtered = append(filtered, val.(string))
	}
	return filtered
}

func Contains(s primitive.A, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
