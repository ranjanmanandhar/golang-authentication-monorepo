package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"gitlab-server.wlink.com.np/nettv/nettv-auth/lib/util"
)

// Headers : Extra headers to be passed to the api
type Headers struct {
	Key, Value string
}

type RequestOptions struct {
	Timeout *time.Duration
}

// Request : Request Options
type Request struct {
	Method   string
	Payload  interface{}
	Endpoint string
	Timeout  *time.Duration
}

// NETTV : For Base URL and Token
type Nettv struct {
	Hostname string
	Schema   string
	Token    string
}

type nettvservice struct {
	logger log.Logger
	nettv  Nettv
}

type NettvStatus struct {
	status int
}

type NettvService interface {
	NewRequest(r *Request, headers []*Headers) (int, []byte, error)
	GetNettvLoginStatus(r *RequestOptions) (int, error)
}

func NewNettvService(logger log.Logger, nettv Nettv) NettvService {
	return &nettvservice{
		logger: logger,
		nettv:  nettv,
	}
}

const (
	bearerSchema = "Basic "
	// Authorization constants
	Authorization = "Authorization"
)

func (q *nettvservice) NewRequest(r *Request, headers []*Headers) (int, []byte, error) {

	url := fmt.Sprintf("%s://%s/%s", q.nettv.Schema, q.nettv.Hostname, r.Endpoint)
	var req *http.Request
	var reqerr error
	if r.Payload != nil {
		buf := new(bytes.Buffer)
		level.Info(q.logger).Log("METHOD", "NewRequest", "msg", "new request stated", "payload", r.Payload)
		json.NewEncoder(buf).Encode((r.Payload))
		req, reqerr = http.NewRequest(r.Method, url, buf)
	} else {
		req, reqerr = http.NewRequest(r.Method, url, nil)
	}
	req.SetBasicAuth("www3", "mmm321")

	if len(headers) > 0 {
		for _, h := range headers {
			req.Header.Add(h.Key, h.Value)
		}
	}
	if reqerr != nil {
		return 500, nil, reqerr
	}

	timeout := 5 * time.Second
	if r.Timeout != nil {
		timeout = *r.Timeout
	}

	client := &http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 500, nil, reqerr
	}
	respData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return resp.StatusCode, nil, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, respData, nil
}

func (n *nettvservice) GetNettvLoginStatus(r *RequestOptions) (int, error) {
	level.Info(n.logger).Log("METHOD", "GetNettvLoginStatus", "msg", "GetNettvLoginStatus stated")

	var num int
	req := &Request{
		Endpoint: "w3gateway/ibill/verify_me/",
		Method:   http.MethodPost,
	}

	if r.Timeout != nil {
		req.Timeout = r.Timeout
	}
	code, data, err := n.NewRequest(req, nil)

	if code != 200 {
		level.Error(n.logger).Log("METHOD", "GetNettvLoginStatus", "msg", data)
		return 0, util.TranslateErrorCode(code)
	}

	if err != nil {
		level.Error(n.logger).Log("METHOD", "GetNettvLoginStatus", "msg", err)
		return 0, err
	}

	json.Unmarshal(data, &num)
	if num == 0 {
		level.Error(n.logger).Log("METHOD", "GetNettvLoginStatus", "msg", err)
		return 0, err
	}
	level.Info(n.logger).Log("METHOD", "GetNettvLoginStatus", "msg", "GetNettvLoginStatus completed", "response", num)
	return num, nil
}
