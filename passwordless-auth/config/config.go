package config

//DBConfig
type DBConfig struct {
	Username            string `mapstructure:"DB_USERNAME"`
	Password            string `mapstructure:"DB_PASSWORD"`
	Host                string `mapstructure:"DB_HOST"`
	Port                int    `mapstructure:"DB_PORT"`
	ServiceName         string `mapstructure:"DB_DATABASE"`
	EserviceUsername    string `mapstructure:"ESERVICE_DB_USERNAME"`
	EservicePassword    string `mapstructure:"ESERVICE_DB_PASSWORD"`
	EserviceHost        string `mapstructure:"ESERVICE_DB_HOST"`
	EservicePort        int    `mapstructure:"ESERVICE_DB_PORT"`
	EserviceServiceName string `mapstructure:"ESERVICE_DB_DATABASE"`
	EbillUsername       string `mapstructure:"EBILL_DB_USERNAME"`
	EbillPassword       string `mapstructure:"EBILL_DB_PASSWORD"`
	EbillHost           string `mapstructure:"EBILL_DB_HOST"`
	EbillPort           int    `mapstructure:"EBILL_DB_PORT"`
	EbillServiceName    string `mapstructure:"EBILL_DB_DATABASE"`
}
