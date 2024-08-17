package services

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/mintance/go-uniqid"
)

const AuthorizationTokenVersion = "0.2.1"

type Auth interface {
	MacIpCheck(string, string) map[string]interface{}
	CreateNettvMigrate(string)
}

type passwordlessAuth struct {
	eservicedb *sql.DB
	ebilldb    *sql.DB
	logger     log.Logger
}

func AuthService(eservicedb *sql.DB, ebilldb *sql.DB, logger log.Logger) Auth {
	return &passwordlessAuth{
		eservicedb: eservicedb,
		ebilldb:    ebilldb,
		logger:     log.With(logger, "AuthService", "MacIPCheck"),
	}
}

func (s *passwordlessAuth) MacIpCheck(mac_id string, stb_ip_address string) map[string]interface{} {
	statement := `SELECT ce.user_name, ce.pay_plan, w.disable, w.skipipcheck, ce.ipoe_mac FROM customer_ebill ce, ebill.wlinkiptv_usermap w
					WHERE ce.user_name = w.username AND w.stb_box_id =:id`
	row := s.eservicedb.QueryRow(statement, mac_id)
	var user_name string
	var pay_plan string
	var disable string
	var skipipcheck string
	var ipoe_mac string
	errScan := row.Scan(&user_name, &pay_plan, &disable, &skipipcheck, &ipoe_mac)
	if errScan != nil || disable == "Y" {
		if errScan != nil {
			level.Error(s.logger).Log("Error scaning row", errScan)
		} else {
			level.Error(s.logger).Log("Disabled User", fmt.Sprintf("Username: %s", user_name))
		}
		return map[string]interface{}{
			"radius_ip_address":      nil,
			"stb_request_ip_address": stb_ip_address,
			"match_ip_address":       false,
			"error": map[string]interface{}{
				"code": 404,
				"msg":  "This device is not registered with any of the account.",
			},
			"authorizer": "worldlink",
		}
	}

	fmt.Printf("Username: %s\n Payplan: %s\n Disable: %s \n Skipipcheck: %s \n Ipoemac: %s\n", user_name, pay_plan, disable, skipipcheck, ipoe_mac)
	var frameAddress string
	if ipoe_mac == "" {
		frameAddress = GetFrameAddress(s.eservicedb, ipoe_mac)
	} else {
		frameAddress = GetFrameAddress(s.eservicedb, user_name)
	}
	if frameAddress == "" {
		return map[string]interface{}{
			"user_name":              user_name,
			"pay_plan":               pay_plan,
			"disable":                disable,
			"skipipcheck":            skipipcheck,
			"ipoe_mac":               ipoe_mac,
			"radius_ip_address":      false,
			"stb_request_ip_address": stb_ip_address,
			"match_ip_address":       false,
			"error":                  false,
			"authorizer":             "worldlink",
		}
	} else {
		return map[string]interface{}{
			"user_name":              user_name,
			"pay_plan":               pay_plan,
			"disable":                disable,
			"skipipcheck":            skipipcheck,
			"ipoe_mac":               ipoe_mac,
			"radius_ip_address":      frameAddress,
			"stb_request_ip_address": stb_ip_address,
			"match_ip_address":       stb_ip_address == frameAddress,
			"error":                  false,
			"authorizer":             "worldlink",
		}
	}
}

func GetFrameAddress(eservicedb *sql.DB, ipoeOrUsername string) string {
	statement := `SELECT
						framedipaddress
					FROM
						(
							SELECT
								framedipaddress
							FROM
								fradius.radacct
							WHERE
								username LIKE concat(:username, '%')
								OR username LIKE concat(:username, '%')
								AND framedipaddress IS NOT NULL
								AND acctstoptime IS NULL
							ORDER BY
								radacctid DESC
						)
					WHERE
						ROWNUM = 1`
	row := eservicedb.QueryRow(statement, ipoeOrUsername)
	var frameAddress string
	errScan := row.Scan(&frameAddress)
	if errScan != nil {
		fmt.Println(errScan)
		return ""
	}
	return frameAddress
}

func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-FORWARDED-FOR")
	if ip == "" {
		ip = r.Header.Get("X-REAL-IP")
		if ip == "" {
			ip = r.RemoteAddr
		}
	}
	return ip
}

func FetchToken(r *http.Request) (string, string) {
	appSecret := r.Header.Get("APPSECRET")
	appID := r.Header.Get("APPID")

	return appSecret, appID
}

func GetAuthorizationToken(payload map[string]interface{}, appSecret string) string {
	id := uniqid.New(uniqid.Params{})
	header := map[string]interface{}{
		"version": AuthorizationTokenVersion,
		"id":      id,
	}
	jsonPayload, _ := json.Marshal(payload)
	jsonHeader, _ := json.Marshal(header)

	signature := md5.Sum([]byte(string(jsonHeader) + string(jsonPayload) + appSecret))

	authToken := map[string]interface{}{
		"header":    header,
		"payload":   payload,
		"signature": hex.EncodeToString(signature[:]),
	}

	jsonToken, _ := json.Marshal(authToken)

	token := base64.StdEncoding.EncodeToString([]byte(string(jsonToken)))

	return token
}

func (s *passwordlessAuth) CreateNettvMigrate(username string) {
	// extra := nil
	// statement := fmt.Sprintf(`insert into NETTV_MIGRATE_TEMP (id,username,extra,status) values(nettv_migrate_temp_seq.nextval, %s,null,1)`, username)
	statement := `select username from NETTV_MIGRATE_TEMP where username='pradimna_home' AND extra is null`
	row, err := s.ebilldb.Query(statement)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(row)
	// // fmt.Println(row)
	// var frameAddress string
	// for row.Next() {
	// 	errScan := row.Scan(&frameAddress)
	// 	fmt.Println(frameAddress)
	// 	if errScan != nil {
	// 		fmt.Println(errScan)
	// 		// return ""
	// 	}
	// }

	// return frameAddress
}
