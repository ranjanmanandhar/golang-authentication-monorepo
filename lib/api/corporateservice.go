package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type corporateservice struct {
	logger log.Logger
}

type ServiceCorporate interface {
	GetCorporateUserDetail(username string) map[string]interface{}
	GetCustomerCodeDetail(customerCode string) map[string]interface{}
}

func CorporateService(logger log.Logger) ServiceCorporate {
	return corporateservice{
		logger: logger,
	}
}

const corporateUrl = "https://corporate-service.wlink.com.np"

func (s corporateservice) GetCorporateUserDetail(username string) map[string]interface{} {
	url := fmt.Sprintf("%s/v2/cr/circuitId?circuit_id=%s", corporateUrl, username)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJuZXR0di1zZXJ2aWNlIiwiaXNzIjoibG9naW4ud2xpbmsuY29tLm5wIiwiaWF0IjoxNjczNDE4NTc0LCJuYmYiOjE2NzM0MTg1NzQsImp0aSI6IjNLN2ZVMmFFeEo5VWJ0elgifQ.uW8FbLUo0FKOwZHlDFwM5jORrKD5qgpRhwsHm_WbHAw")
	if err != nil {
		level.Info(s.logger).Log(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		level.Error(s.logger).Log(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			level.Info(s.logger).Log(err)
		}
		var data map[string]interface{}
		json.Unmarshal(body, &data)
		level.Info(s.logger).Log("msg", "get user data of corporate success "+username, "RESPONSE", string(body))
		return data
	}
	errResp, _ := ioutil.ReadAll(resp.Body)
	level.Error(s.logger).Log("msg", "get user status error response ", "RESPONSE", errResp)
	if err != nil {
		level.Info(s.logger).Log(err)
	}
	return nil
}

func (s corporateservice) GetCustomerCodeDetail(customerCode string) map[string]interface{} {
	url := fmt.Sprintf("%s/v2/customer-code/%s", corporateUrl, customerCode)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJuZXR0di1zZXJ2aWNlIiwiaXNzIjoibG9naW4ud2xpbmsuY29tLm5wIiwiaWF0IjoxNjczNDE4NTc0LCJuYmYiOjE2NzM0MTg1NzQsImp0aSI6IjNLN2ZVMmFFeEo5VWJ0elgifQ.uW8FbLUo0FKOwZHlDFwM5jORrKD5qgpRhwsHm_WbHAw")
	if err != nil {
		level.Info(s.logger).Log(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		level.Error(s.logger).Log(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			level.Info(s.logger).Log(err)
		}
		var data map[string]interface{}
		json.Unmarshal(body, &data)
		level.Info(s.logger).Log("msg", "get customer code data success "+customerCode, "RESPONSE", string(body))
		return data
	}
	errResp, _ := ioutil.ReadAll(resp.Body)
	level.Error(s.logger).Log("msg", "get customer code data error response ", "RESPONSE", errResp)
	if err != nil {
		level.Info(s.logger).Log(err)
	}
	return nil
}
