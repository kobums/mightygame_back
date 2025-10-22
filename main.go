package main

import (
	"mighty/config"
	"mighty/global"
	"mighty/models"
	"mighty/router"
	"mighty/services"
	"net/http"
	"time"

	"encoding/gob"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&models.User{})

	//rand.Seed(time.Now().UnixNano())
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	log.Printf("mighty Version=" + config.Version + " Build=" + config.Build)

	log.Info("Server Start")

	// DB connection is not required for in-memory room system
	// models.InitCache()

	go services.Cron()

	Http()
}

func Http() {
	r := gin.Default()

	c := cors.DefaultConfig()
	c.AllowOrigins = []string{"http://localhost:8080"}

	r.Use(cors.New(c))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.Static("/assets", "./assets")
	r.Static("/webdata", "./webdata")

	r.GET("/", func(c *gin.Context) {
		content := global.ReadFile("./views/index.html")

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, content)
	})

	router.SetRouter(r)

	s := &http.Server{
		Addr:           ":" + config.Port,
		Handler:        r,
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
