package main

import (
	"crypto/rand"
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

	// 하드코딩된 쿠키 서명 키("secret")를 쓰면 누구나 세션을 위조할 수 있다.
	// SESSION_SECRET 이 없으면 기동마다 랜덤 키를 생성한다(재시작 시 세션 만료).
	sessionKey := []byte(config.SessionSecret)
	if len(sessionKey) == 0 {
		sessionKey = make([]byte, 32)
		if _, err := rand.Read(sessionKey); err != nil {
			log.Fatal(err)
		}
		log.Warn("SESSION_SECRET not set: using a random session key")
	}
	store := cookie.NewStore(sessionKey)
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
