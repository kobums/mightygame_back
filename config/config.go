package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	Database         string
	Owner            string
	ConnectionString string
	Port             string
	TempPath         string

	ServiceUrl string

	UploadPath string

	SmsUser   string
	SmsKey    string
	SmsSender string

	AdminEmail string

	Version string
	Build   string
	DEBUG   uint64

	Domain string
)

func init() {
	UploadPath = "webdata"
	Database = "mysql"
	Port = "9004"
	Domain = "localhost"

	DEBUG = 0
	if os.Getenv("GIN_MODE") == "release" {
		DEBUG = 0
	}
	if DEBUG > 0 {
		fmt.Printf("Debug: MODE=true, flag=%+v \n", DEBUG)
	}

	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if value := viper.Get("connectionString"); value != nil {
		ConnectionString = value.(string)
	}

	if value := viper.Get("uploadPath"); value != nil {
		UploadPath = value.(string)
	}

	if value := viper.Get("port"); value != nil {
		Port = value.(string)
	}

	if value := viper.Get("domain"); value != nil {
		Domain = value.(string)
	}
}
