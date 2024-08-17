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

type corporatenettvservice struct {
	logger log.Logger
}

type CorporateNettv interface {
	GetCorporateNettv(username string) []string
	CheckForNewSystem(username string) bool
	GetCorporateStbBoxes(username string) ([]string, string, bool)
}

func CorporateNettvService(logger log.Logger) CorporateNettv {
	return corporatenettvservice{
		logger: logger,
	}
}

// const cndUrl = "https://cnd.wlink.com.np/api"
const cndUrl = "http://10.12.7.99:8080"

func (s corporatenettvservice) GetCorporateNettv(username string) []string {
	level.Info(s.logger).Log("msg", "Get Stbs of Corporate User for "+username)
	url := fmt.Sprintf("%s/users/%s/stbs", cndUrl, username)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Exchange-Username", "esupport")
	if err != nil {
		level.Info(s.logger).Log(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		level.Error(s.logger).Log("error", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			level.Info(s.logger).Log(err)
		}
		var data []interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("err", err)
		}
		stbs := mapStbsForNewSystem(data)
		fmt.Println("stbs", stbs)
		if len(data) <= 0 {
			return stbs
		}

		level.Info(s.logger).Log("msg", "get user data of corporate success "+username, "RESPONSE", stbs)
		return stbs
	}
	errResp, _ := ioutil.ReadAll(resp.Body)
	level.Error(s.logger).Log("msg", "get user status error response ", "RESPONSE", errResp)
	if err != nil {
		level.Info(s.logger).Log(err)
	}
	return nil

}

func (s corporatenettvservice) CheckForNewSystem(username string) bool {
	url := fmt.Sprintf("%s/users/%s/check", cndUrl, username)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Exchange-Username", "esupport")
	if err != nil {
		level.Info(s.logger).Log(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		level.Error(s.logger).Log("error", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			level.Info(s.logger).Log(err)
		}
		var data map[string]interface{}
		json.Unmarshal(body, &data)
		level.Info(s.logger).Log("msg", "check system of corporate success "+username, "RESPONSE", string(body))
		sys := data["data"].(map[string]interface{})
		return sys["user_exists"].(bool)
	}
	level.Info(s.logger).Log("msg", "check system of corporate failed "+username, "RESPONSE", resp)
	return false
}

func (s corporatenettvservice) GetCorporateStbBoxes(username string) ([]string, string, bool) {
	checkSystem := s.CheckForNewSystem(username)
	var stbs []string
	if checkSystem {
		stbs = s.GetCorporateNettv(username)
		return stbs, "Y", true
	}
	return stbs, "", false
}

func mapStbsForNewSystem(data []interface{}) []string {
	var stbs []string
	for _, val := range data {
		newVal := val.(map[string]interface{})
		currentStb := fmt.Sprintf("%s", newVal["mac_address"])
		stbs = append(stbs, currentStb)
	}
	return stbs
}
