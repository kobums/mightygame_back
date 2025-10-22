package router

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func JwtAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")

		if user != nil {
			c.Next()
			return
		}

		var str string
		str = "Not Auth"
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte(str))
		c.Abort()
	}
}
