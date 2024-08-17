package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type nettvService struct {
	logger log.Logger
}

type NettvService interface {
	GetCustomerNettvDetail(username string) []map[string]interface{}
}

const NettvServiceUrl = "https://api-nettv.wlink.com.np"

func NewNettvService(logger log.Logger) NettvService {
	return nettvService{
		logger: logger,
	}
}

func (n nettvService) GetCustomerNettvDetail(username string) []map[string]interface{} {
	level.Info(n.logger).Log("msg", "Get customer nettv detail started ", "payload", "data")
	url := fmt.Sprintf("%s/v1/users/%s/detail", NettvServiceUrl, username)
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		level.Info(n.logger).Log(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		level.Error(n.logger).Log(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			level.Error(n.logger).Log(err)
		}
		var data []map[string]interface{}
		json.Unmarshal(body, &data)
		level.Info(n.logger).Log("msg", "get customer nettv detail success", "RESPONSE", string(body))
		return data
	}
	errResp, _ := ioutil.ReadAll(resp.Body)
	level.Error(n.logger).Log("msg", "get customer nettv detail error response ", "RESPONSE", errResp)
	if err != nil {
		level.Error(n.logger).Log(err)
	}
	return nil
}
