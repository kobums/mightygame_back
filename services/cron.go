package services

import (
	"github.com/robfig/cron"

	log "github.com/sirupsen/logrus"
)

func Cron() {
	log.Println("Cron Service Start")

	c := cron.New()

	//c.AddFunc("0 * * * * *", ResetThrottling)

	//c.AddFunc("0 * * * * *", SendEmail)
	//c.AddFunc("0 * * * * *", SendSMS)

	c.Start()
}
