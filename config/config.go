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

	// 세션 쿠키 서명 키. 미설정 시 기동마다 랜덤 생성(main.go).
	SessionSecret string
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
	// config.json 은 DB 접속 정보를 담으므로 커밋하지 않는다(config.json.example 참고).
	// 파일이 없으면 환경변수(DATABASE_URL)만으로도 기동할 수 있다.
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("config file not loaded: %s \n", err)
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

	// 환경변수가 있으면 설정 파일 값보다 우선한다.
	if value := os.Getenv("DATABASE_URL"); value != "" {
		ConnectionString = value
	}
	if value := os.Getenv("SESSION_SECRET"); value != "" {
		SessionSecret = value
	}
}
