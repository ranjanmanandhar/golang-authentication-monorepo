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

type servicewlink struct {
	logger log.Logger
}

type ServiceWlink interface {
	GetUserStatus(username string) map[string]interface{}
	GetCustomerAccountDetail(plan_category_id string) map[string]interface{}
	GetCustomerInfo(username string) map[string]interface{}
	GetCustomerDetailsWithPassword(username string) map[string]interface{}
}

func WlinkService(logger log.Logger) ServiceWlink {
	return servicewlink{
		logger: logger,
	}
}

const servicewlinkUrl = "https://services.wlink.com.np"

func Greet(audience string) string {
	return fmt.Sprintf("H, %s!", audience)
}

func (s servicewlink) GetUserStatus(username string) map[string]interface{} {
	level.Info(s.logger).Log("msg", "get user status started "+username, "payload", "data")
	url := fmt.Sprintf("%s/customers/customers/%s/status", servicewlinkUrl, username)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("ebill", "RPx9arjSfJ5BVDpN")
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
		level.Info(s.logger).Log("msg", "get user status success "+username, "RESPONSE", string(body))
		return data
	}
	errResp, _ := ioutil.ReadAll(resp.Body)
	level.Error(s.logger).Log("msg", "get user status error response ", "RESPONSE", errResp)
	if err != nil {
		level.Info(s.logger).Log(err)
	}
	return nil
}

func (s servicewlink) GetCustomerAccountDetail(plan_category_id string) map[string]interface{} {
	level.Info(s.logger).Log("msg", "Get user bandwidth started ", "payload", "data")
	url := fmt.Sprintf("%s/plans/plancategories/%s", servicewlinkUrl, plan_category_id)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("ebill", "RPx9arjSfJ5BVDpN")
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
			level.Error(s.logger).Log(err)
		}
		var data map[string]interface{}
		json.Unmarshal(body, &data)
		level.Info(s.logger).Log("msg", "get user bandwidth success", "RESPONSE", string(body))
		return data
	}
	errResp, _ := ioutil.ReadAll(resp.Body)
	level.Error(s.logger).Log("msg", "get user bandwidth error response ", "RESPONSE", errResp)
	if err != nil {
		level.Error(s.logger).Log(err)
	}
	return nil
}

func (s servicewlink) GetCustomerInfo(username string) map[string]interface{} {
	level.Info(s.logger).Log("msg", "Get customer info started ", "payload", "data")
	url := fmt.Sprintf("%s/customers/customerinfos/%s", servicewlinkUrl, username)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("ebill", "RPx9arjSfJ5BVDpN")
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
			level.Error(s.logger).Log(err)
		}
		var data map[string]interface{}
		json.Unmarshal(body, &data)
		level.Info(s.logger).Log("msg", "get customer info success", "RESPONSE", string(body))
		return data
	}
	errResp, _ := ioutil.ReadAll(resp.Body)
	level.Error(s.logger).Log("msg", "get customer info error response ", "RESPONSE", errResp)
	if err != nil {
		level.Error(s.logger).Log(err)
	}
	return nil
}

func (s servicewlink) GetCustomerDetailsWithPassword(username string) map[string]interface{} {
	level.Info(s.logger).Log("msg", "Get customer details with password started ", "payload", "data")
	url := fmt.Sprintf("%s/customers/customers/%s", servicewlinkUrl, username)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth("ebill", "RPx9arjSfJ5BVDpN")
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
			level.Error(s.logger).Log(err)
		}
		var data map[string]interface{}
		json.Unmarshal(body, &data)
		level.Info(s.logger).Log("msg", "Get customer details with password success", "RESPONSE", string(body))
		return data
	}
	errResp, _ := ioutil.ReadAll(resp.Body)
	level.Error(s.logger).Log("msg", "Get customer details with password error response ", "RESPONSE", errResp)
	if err != nil {
		level.Error(s.logger).Log(err)
	}
	return nil
}
