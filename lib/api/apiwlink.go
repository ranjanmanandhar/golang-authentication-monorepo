package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/types"
)

type apiwlink struct {
	logger log.Logger
}

type ApiWlink interface {
	AdAuthenticate(username string, password string) types.OldAdminAuthResponse
}

func ApiWlinkService(logger log.Logger) ApiWlink {
	return apiwlink{
		logger: logger,
	}
}

const apiwlinkUrl = "http://api.wlink.com.np/core/ad/authenticate"

func (a apiwlink) AdAuthenticate(username string, password string) types.OldAdminAuthResponse {
	level.Info(a.logger).Log("METHOD", "AdAuthenticate", "msg", "AdAuthenticate started ")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	form := url.Values{}
	form.Add("username", username)
	form.Add("password", password)

	req, err := http.NewRequest("POST", apiwlinkUrl, strings.NewReader(form.Encode()))
	if err != nil {
		level.Error(a.logger).Log("METHOD", "AdAuthenticate", "msg", err)
	}

	req.SetBasicAuth("da39a3ee5e6b4b0d3255bfef95601890afd80709", "x")

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		level.Error(a.logger).Log("METHOD", "AdAuthenticate", "msg", err)
	}
	defer resp.Body.Close()

	var data types.OldAdminAuthResponse

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			level.Info(a.logger).Log(err)
		}

		json.Unmarshal(body, &data)
		level.Info(a.logger).Log("METHOD", "AdAuthenticate", "msg", "Adl authenticate success", "RESPONSE", string(body))
		return data
	}

	return data
}
