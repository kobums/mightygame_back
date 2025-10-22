package api

import (
	"mighty/controllers"
)

type Login struct {
	controllers.Controller
}

func (c *Login) Login(loginid string, passwd string) {
}
