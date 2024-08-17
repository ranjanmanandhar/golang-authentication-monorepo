package types

import (
	"database/sql"
	"time"
)

type CustomerInfo struct {
	Username            string `json:"USERNAME" bson:"username,omitempy"`
	FullName            string `json:"NAME" bson:"fullName,omitempy"`
	Address             string `json:"ADDRESS" bson:"address,omitempy"`
	PHONE1              string `json:"PHONE1" bson:"phone1,omitempy"`
	PHONE2              string `json:"PHONE2" bson:"phone2,omitempy"`
	PrimaryMobileNumber string `json:"MOBILE" bson:"primaryMobileNumber,omitempy"`
	PrimaryEmailAddress string `json:"EMAIL" bson:"primaryEmailAddress,omitempy"`
}

type CustomerInfoData struct {
	Username     string       `json:"USERNAME" bson:"username,omitempty"`
	CustomerInfo CustomerInfo `json:"CUST_INFO" bson:"customerInfo,omitempty"`
	Updated_At   time.Time    `json:"UPDATED_AT" bson:"updated_at,omitempty"`
	Data_type    string       `json:"DATA_TYPE" bson:"data_type,omitempty"`
}

type AccountStatus struct {
	Username            string `json:"USERNAME" bson:"username,omitempty" redis:"USERNAME"`
	PayPlan             string `json:"PAY_PLAN" bson:"payplan" redis:"PAY_PLAN"`
	PlanName            string `json:"PLAN_NAME" bson:"planname,omitempty" redis:"PLAN_NAME"`
	ExpiryDate          string `json:"EXPIRY_DATE" bson:"expirydate,omitempty" redis:"EXPIRY_DATE"`
	DaysRemaining       string `json:"DAYS_LEFT" bson:"daysremaining" redis:"DAYS_LEFT"`
	MinuteLeft          string `json:"MINUTE_LEFT" bson:"minuteleft,omitempty" redis:"MINUTE_LEFT"`
	VolumeLeft          string `json:"VOLUME_LEFT" bson:"volumeleft,omitempty" redis:"VOLUME_LEFT"`
	LastOnline          string `json:"LAST_ONLINE" bson:"lastonline,omitempty" redis:"LAST_ONLINE"`
	SubscribedBandwidth string `json:"SUBSCRIBED_BANDWIDTH" bson:"subscribedbandwidth,omitempty" redis:"SUBSCRIBED_BANDWIDTH"`
	CurrentBandwidth    string `json:"CURRENT_BANDWIDTH" bson:"currentbandwidth,omitempty" redis:"CURRENT_BANDWIDTH"`
	FallbackSpeed       string `json:"FALLBACK_SPEED" bson:"fallbackspeed,omitempty" redis:"FALLBACK_SPEED"`
	Balance             string `json:"BALANCE" bson:"balance,omitempty" redis:"BALANCE"`
	GraceStatus         string `json:"GRACE_STATUS" bson:"gracestatus" redis:"GRACE_STATUS"`
	Disable             string `json:"DISABLE" bson:"disable" redis:"DISABLE"`
	AccountStatus       string `json:"ACCOUNT_STATUS" bson:"accountstatus,omitempty" redis:"ACCOUNT_STATUS"`
	IpoeMac             string `json:"IPOE_MAC" bson:"ipoeMac,omitempty" redis:"IPOE_MAC"`
	Password            string `json:"PASSWORD" bson:"password" redis:"PASSWORD"`
}

type AccountStatusData struct {
	Username      string        `json:"USERNAME" bson:"username,omitempty"`
	AccountStatus AccountStatus `json:"ACCOUNT_INFO" bson:"accountInfo,omitempty"`
	Updated_At    time.Time     `json:"UPDATED_AT" bson:"updated_at,omitempty"`
	Data_type     string        `json:"DATA_TYPE" bson:"data_type,omitempty"`
}

type Data struct {
	Username          string                 `json:"USERNAME" bson:"username,omitempty" redis:"USERNAME"`
	CustomerInfo      CustomerInfo           `json:"CUST_INFO" bson:"customerInfo,omitempty" redis:"CUST_INFO"`
	AccountStatus     AccountStatus          `json:"ACCOUNT_INFO" bson:"accountInfo,omitempty" redis:"ACCOUNT_INFO"`
	STB_BOX_ID        []string               `json:"-" bson:"nettvsettopboxes,omitempty" redis:"STB_BOX_ID"`
	MAPPED_STB_BOX_ID map[string]interface{} `json:"STB_BOX_ID" bson:"mappednettvsettopboxes,omitempty" redis:"STB_BOX_ID"`
	Session_Count     string                 `json:"SESSION_COUNT" bson:"session_count,omitempty" redis:"SESSION_COUNT"`
	Skipipcheck       string                 `json:"SKIPIPCHECK" bson:"skipipcheck,omitempty" redis:"SKIPIPCHECK"`
	New_Nettv_Map     string                 `json:"NEW_NETTV_MAP" bson:"new_nettv_map,omitempty" redis:"NEW_NETTV_MAP"`
	Is_Nettv_Disabled string                 `json:"IS_NETTV_DISABLED" bson:"is_nettv_disabled,omitempty" redis:"IS_NETTV_DISABLED"`
	Updated_At        time.Time              `json:"UPDATED_AT" bson:"updated_at,omitempty" redis:"UPDATED_AT"`
	Data_Type         string                 `json:"DATA_TYPE" bson:"data_type,omitempty" redis:"DATA_TYPE"`
}

type NettvPayload struct {
	Username   string `json:"username,omitempty"`
	StbBoxIds  string `json:"stb_box_ids,omitempty"`
	NewNetvMap string `json:"new_netv_map,omitempty"`
	DataType   string `json:"data_type,omitempty"`
	Type       string `json:"type,omitempty"`
}

type RetailNettvData struct {
	Username          string    `json:"username,omitempty"`
	STB_BOX_ID        []string  `json:"STB_BOX_ID" bson:"nettvsettopboxes"`
	New_Nettv_Map     string    `json:"NEW_NETTV_MAP" bson:"new_nettv_map,omitempty"`
	Session_Count     string    `json:"SESSION_COUNT" bson:"session_count,omitempty"`
	Skip_Ip_Check     string    `json:"SKIPIPCHECK" bson:"skipipcheck,omitempty"`
	Is_Nettv_Disabled string    `json:"IS_NETTV_DISABLED" bson:"is_nettv_disabled,omitempty"`
	Updated_At        time.Time `json:"UPDATED_AT" bson:"updated_at,omitempty"`
	Data_type         string    `json:"DATA_TYPE" bson:"data_type,omitempty"`
}

type CorporateNettvPayload struct {
	Username        string   `json:"username,omitempty"`
	CompanyName     string   `json:"company_name,omitempty"`
	StbBoxIds       []string `json:"stb_box_id,omitempty"`
	Address         string   `json:"address,omitempty"`
	Email           string   `json:"email,omitempty"`
	Phone           string   `json:"phone,omitempty"`
	Mobile          string   `json:"mobile_no,omitempty"`
	DataType        string   `json:"data_type,omitempty"`
	Type            string   `json:"type,omitempty"`
	CustomerCode    string   `json:"customer_code"`
	Password        string   `json:"password"`
	AccountStatus   string   `json:"ACCOUNT_STATUS" bson:"account_status"`
	CircuitIdStatus string   `json:"CIRCUIT_ID_STATUS" bson:"circuit_id_status"`
}

type CircuitData struct {
	Username        string                `json:"USERNAME,omitempty" bson:"username"`
	CustomerInfo    CorporateCustomerInfo `json:"CUSTOMER_INFO" bson:"customer_info"`
	CustomerCode    string                `json:"CUSTOMER_CODE" bson:"customer_code"`
	AccountStatus   string                `json:"ACCOUNT_STATUS" bson:"account_status"`
	CircuitIdStatus string                `json:"CIRCUIT_ID_STATUS" bson:"circuit_id_status"`
	Updated_At      time.Time             `json:"UPDATED_AT" bson:"updated_at,omitempty"`
	Data_type       string                `json:"DATA_TYPE" bson:"data_type,omitempty"`
}

type CorporateCustomerInfo struct {
	Username string `json:"USERNAME"`
	Customer string `json:"CUSTOMER"`
	Address  string `json:"ADDRESS"`
	Email    string `json:"EMAIL"`
	Phone1   string `json:"PHONE1"`
	Mobile   string `json:"MOBILE"`
}

type CorporateData struct {
	Username        string                `json:"USERNAME,omitempty" bson:"username,omitempty"`
	CustomerCode    string                `json:"-" bson:"customer_code,omitempty"`
	CustomerInfo    CorporateCustomerInfo `json:"CUSTOMER_INFO"`
	StbBoxId        []string              `json:"STB_BOX_ID"`
	NewNettvMap     string                `json:"NEW_NETTV_MAP" bson:"new_nettv_map,omitempty"`
	IsNettvDisabled string                `json:"IS_NETTV_DISABLED" bson:"is_nettv_disabled,omitempty"`
	SessionCount    string                `json:"SESSION_COUNT" bson:"session_count,omitempty"`
	Password        string                `json:"-" bson:"password,omitempty"`
	AccountStatus   string                `json:"-" bson:"account_status,omitempty"`
	Updated_At      time.Time             `json:"-" bson:"updated_at,omitempty"`
	Data_type       string                `json:"-" bson:"data_type,omitempty"`
}
type CorporateNettvData struct {
	Username string   `json:"username"`
	StbBoxId []string `json:"stb_box_id"`
}
type NewCspayload struct {
	Username               string `json:"user_name"`
	Machine_Name           string `json:"machine_name"`
	Subscription_No        string `json:"subscription_no"`
	Account_Type           string `json:"account_type"`
	Pay_Plan               string `json:"pay_plan"`
	Mem_Start_Date         string `json:"mem_start_date"`
	Mem_Expire_Date        string `json:"mem_expire_date"`
	Per_Client_Name        string `json:"per_client_name"`
	Per_Address            string `json:"per_address"`
	Per_Cont_Mobile        string `json:"per_cont_mobile"`
	Per_Cont_Email_Primary string `json:"per_cont_email_primary"`
	Installed_By           string `json:"installed_by"`
	Marketed_By            string `json:"marketed_by"`
	Referred_By            string `json:"referred_by"`
	Plan_Category_Id       string `json:"plan_category_id"`
	Disable                string `json:"disable"`
	Per_Ward_No            string `json:"per_ward_no"`
	Per_House_No           string `json:"per_house_no"`
	Latitude               string `json:"latitude"`
	Longitude              string `json:"longitude"`
	Tt_Ni_Ticketid         string `json:"tt_ni_ticketid"`
}

type CorporateCustomerCodeData struct {
	CustomerCode  string    `json:"customer_code" bson:"customer_code"`
	Password      string    `json:"password" bson:"password"`
	AccountStatus string    `json:"account_status" bson:"account_status"`
	Updated_At    time.Time `json:"UPDATED_AT" bson:"updated_at,omitempty"`
}

type CorporateReturnData struct {
	Username      string `json:"USERNAME"`
	Customer      string `json:"CUSTOMER"`
	Address       string `json:"ADDRESS"`
	Email         string `json:"EMAIL"`
	Phone1        string `json:"PHONE1"`
	Mobile        string `json:"MOBILE"`
	DISABLE       string `json:"DISABLE"`
	NEW_NETTV_MAP string `json:"NEW_NETTV_MAP"`
}

type CustomerCodeFallbackResponse struct {
	Code    string                             `json:"code" mapstructure:"code"`
	Data    []CustomerCodeFallbackResponseData `json:"data" mapstructure:"data"`
	Message string                             `json:"message" mapstructure:"message"`
}

type CustomerCodeFallbackResponseData struct {
	Name     string `json:"name"`
	Id       string `json:"id"`
	Code     string `json:"code"`
	Password string `json:"password"`
}

type CorporateCustFallbackResponseData struct {
	Username      string `json:"circuit_id"`
	Customer      string `json:"customer_name"`
	Address       string `json:"site_address"`
	Email         string `json:"branch_email"`
	Phone1        string `json:"branch_phone"`
	Mobile        string `json:"branch_contact_number"`
	Password      string `json:"password"`
	CustomerCode  string `json:"customer_code"`
	AccountStatus string `json:"account_status"`
}

type CorporateCustFallbackResponse struct {
	Code    string                              `json:"code" mapstructure:"code"`
	Data    []CorporateCustFallbackResponseData `json:"data" mapstructure:"data"`
	Message string                              `json:"message" mapstructure:"message"`
}

type RedisHealthCheckResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type OldSystemMsg struct {
	Message interface{} `json:"msg"`
}
type OldSystemResponse struct {
	Status int8         `json:"status"`
	Data   OldSystemMsg `json:"data"`
}
type OldAdminAuthResponseData struct {
	Samaccountname string `json:"samaccountname"`
	Displayname    string `json:"displayname"`
	Mail           string `json:"mail"`
	Mobile         string `json:"mobile"`
	Department     string `json:"department"`
}
type OldAdminAuthResponse struct {
	StatusCode int                      `json:"status_code"`
	Data       OldAdminAuthResponseData `json:"data"`
}
type OldAuthPayload struct {
	Username   string `json:"username,omitempty" validate:"required" form:"username"`
	Password   string `json:"password,omitempty" validate:"required_without:Token" form:"password"`
	Mac        string `json:"mac,omitempty" validate:"required" form:"mac"`
	Token      string `json:"token,omitempty" validate:"required_without:Password" form:"token"`
	Ip         string `json:"ip,omitempty" form:"ip"`
	DeviceType string `json:"device_type,omitempty" form:"device_type"`
}

type AdlAuthenticationPayload struct {
	Username string `json:"username,omitempty" validate:"required" form:"username"`
	Password string `json:"password,omitempty" validate:"required" form:"password"`
}

type MacwisePayload struct {
	Mac string `json:"mac,omitempty" validate:"required" form:"mac"`
}
type OldSysUser struct {
	Username    string         `json:"USERNAME" bson:"username,omitempty"`
	Token       string         `json:"token" bson:"token,omitempty"`
	StbBoxId    string         `json:"stb_box_id" bson:"stb_box_id,omitempty"`
	SkipIpCheck sql.NullString `json:"SKIPIPCHECK" bson:"skipipcheck,omitempty"`
	Updated_At  time.Time      `json:"UPDATED_AT" bson:"updated_at,omitempty" redis:"UPDATED_AT"`
}

type OldSysData struct {
	Username          string                 `json:"USERNAME" bson:"username,omitempty" redis:"USERNAME"`
	CustomerInfo      CustomerInfo           `json:"CUST_INFO" bson:"customerInfo,omitempty" redis:"CUST_INFO"`
	AccountStatus     AccountStatus          `json:"ACCOUNT_INFO" bson:"accountInfo,omitempty" redis:"ACCOUNT_INFO"`
	STB_BOX_ID        []string               `json:"-" bson:"nettvsettopboxes,omitempty" redis:"STB_BOX_ID"`
	MAPPED_STB_BOX_ID map[string]interface{} `json:"STB_BOX_ID" bson:"mappednettvsettopboxes,omitempty" redis:"STB_BOX_ID"`
	Session_Count     string                 `json:"SESSION_COUNT" bson:"session_count,omitempty" redis:"SESSION_COUNT"`
	Skipipcheck       string                 `json:"SKIPIPCHECK" bson:"skipipcheck,omitempty" redis:"SKIPIPCHECK"`
	New_Nettv_Map     string                 `json:"NEW_NETTV_MAP" bson:"new_nettv_map,omitempty" redis:"NEW_NETTV_MAP"`
	Is_Nettv_Disabled string                 `json:"IS_NETTV_DISABLED" bson:"is_nettv_disabled,omitempty" redis:"IS_NETTV_DISABLED"`
	Updated_At        time.Time              `json:"UPDATED_AT" bson:"updated_at,omitempty" redis:"UPDATED_AT"`
	Data_Type         string                 `json:"DATA_TYPE" bson:"data_type,omitempty" redis:"DATA_TYPE"`
	Token             string                 `json:"token" bson:"token" redis:"token"`
	IsStbIp           bool                   `json:"isSTBIP" redis:"isSTBIP"`
}
