package main

import (
	"crypto/rand"
	"mighty/config"
	"mighty/global"
	"mighty/models"
	"mighty/router"
	"mighty/services"
	"net/http"
	"os"
	"strings"
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
	// 세션 쿠키를 쓰므로(credentials: 'include') 허용 origin 을 명시하고 AllowCredentials 를 켠다.
	// CORS_ALLOWED_ORIGINS(쉼표 구분)로 운영 origin 을 지정한다. 미설정 시 로컬 개발용 기본값.
	c.AllowOrigins = []string{"http://localhost:8080"}
	if v := os.Getenv("CORS_ALLOWED_ORIGINS"); v != "" {
		var origins []string
		for _, o := range strings.Split(v, ",") {
			if o = strings.TrimSpace(o); o != "" {
				origins = append(origins, o)
			}
		}
		if len(origins) > 0 {
			c.AllowOrigins = origins
		}
	}
	c.AllowCredentials = true

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
	// HttpOnly: JS 에서 세션 쿠키를 읽지 못하게. SameSite=Lax: 타 사이트발 POST 에 쿠키가 실리지 않게(CSRF 완화).
	// HTTPS 로 운영하면 SESSION_COOKIE_SECURE=true 로 Secure 플래그를 켠다.
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   os.Getenv("SESSION_COOKIE_SECURE") == "true",
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("mysession", store))
	// REST 쓰기 엔드포인트(/api/game, /api/gameuser, /api/user)는 로그인 세션 필요 (회원가입 제외).
	r.Use(router.LoginRequiredForWrites())

	r.Static("/assets", "./assets")
	r.Static("/webdata", "./webdata")

	r.GET("/", func(c *gin.Context) {
		content := global.ReadFile("./views/index.html")

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, content)
	})

	// 로그인(POST /api/login/login)·로그아웃·내 정보. router.go 재생성과 무관하게 유지된다.
	router.SetAuthRouter(r)
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
