package types

type User struct {
	AccountInfo  AccountInfo  `json:"ACCOUNT_INFO"`
	Auth         Auth         `json:"AUTH"`
	CustomerInfo CustomerInfo `json:"CUSTOMER_INFO"`
	NettvInfo    NettvInfo    `json:"NETTV_INFO"`
}

type AccountInfo struct {
	AccountStatus       string      `json:"ACCOUNT_STATUS" bson:"account_status"`
	Balance             string      `json:"BALANCE" bson:"balance"`
	CurrentBandwidth    string      `json:"CURRENT_BANDWIDTH" bson:"current_bandwidth"`
	DaysLeft            string      `json:"DAYS_LEFT" bson:"rem_days_left"`
	Disable             string      `json:"DISABLE" bson:"disable"`
	ExpiryDate          string      `json:"EXPIRY_DATE" bson:"expiry_date"`
	FallbackSpeed       string      `json:"FALLBACK_SPEED" bson:"fallback_speed"`
	GraceStatus         string      `json:"GRACE_STATUS" bson:"grace_status"`
	LastOnline          string      `json:"LAST_ONLINE" bson:"online_date"`
	MinuteLeft          interface{} `json:"MINUTE_LEFT" bson:"rem_min_date"`
	PayPlan             string      `json:"PAY_PLAN" bson:"pay_plan"`
	PlanName            string      `json:"PLAN_NAME" bson:"plan_name"`
	SubscribedBandwidth string      `json:"SUBSCRIBED_BANDWIDTH" bson:"subscribed_bandwidth"`
	Username            string      `json:"USERNAME" bson:"username"`
	VolumeLeft          string      `json:"VOLUME_LEFT" bson:"VOLUME_LEFT"`
}

type CustomerInfo struct {
	Address string      `json:"ADDRESS" bson:"ADDRESS"`
	Email   string      `json:"EMAIL" bson:"EMAIL"`
	Mobile  string      `json:"MOBILE" bson:"MOBILE"`
	Name    string      `json:"NAME" bson:"NAME"`
	PHONE1  interface{} `json:"PHONE1" bson:"PHONE1"`
	PHONE2  interface{} `json:"PHONE2" bson:"PHONE2"`
}

type NettvInfo struct {
	Disable         string   `json:"DISABLE"`
	IsNettvDisabled string   `json:"IS_NETTV_DISABLED"`
	NewNettvMap     string   `json:"NEW_NETTV_MAP"`
	SessionCount    string   `json:"SESSION_COUNT"`
	Skipipcheck     bool     `json:"SKIPIPCHECK"`
	Stbs            []string `json:"STBS"`
}

type Auth struct {
	Username string `json:"USERNAME"`
	Password int64  `json:"PASSWORD"`
}

type DefaultAuthenticate struct {
	Authenticate bool        `json:"AUTHENTICATE" bson:"authenticate"`
	Data         AccountInfo `json:"data"`
}

type VerifyClientData struct {
	AccountInfo     AccountInfo       `json:"ACCOUNT_INFO" bson:"ACCOUNT_INFO"`
	CUST_INFO       CustomerInfo      `json:"CUSTOMER_INFO" bson:"CUST_INFO"`
	Disable         string            `json:"DISABLE" bson:"disable"`
	IsNettvDisabled string            `json:"IS_NETTV_DISABLED" bson:"IS_NETTV_DISABLED"`
	NewNettvMap     string            `json:"NEW_NETTV_MAP" bson:"NEW_NETTV_MAP"`
	SessionCount    string            `json:"SESSION_COUNT" bson:"SESSION_COUNT"`
	Skipipcheck     bool              `json:"SKIPIPCHECK"`
	Stbs            map[string]string `json:"STBS" bson:"stbs"`
}

type VerifyClient struct {
	Status int              `json:"status"`
	Data   VerifyClientData `json:"data" bson:"data"`
}
