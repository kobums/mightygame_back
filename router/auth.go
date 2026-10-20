package router

import (
	"net/http"
	"strings"

	"mighty/controllers/api"
	"mighty/models"

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

// 로그인한 세션이 있어야 하는 쓰기(POST/PUT/DELETE/PATCH) REST 엔드포인트.
// /api/room/*, /api/game/state|bid 등 게임 진행 API 는 게스트(로그인 없이 playerName)로 동작하므로 제외한다.
var loginRequiredWritePaths = map[string]bool{
	"/api/game":     true,
	"/api/gameuser": true,
	"/api/user":     true,
}

// LoginRequiredForWrites 는 REST 쓰기 엔드포인트에 로그인 세션을 요구한다.
// router.go 는 buildtool-router 가 재생성하므로, 재생성돼도 유지되도록 이 미들웨어를 main.go 에서 엔진 전체에 건다.
// 회원가입(POST /api/user)은 예외. 본인 계정만 수정/삭제 가능한지는 rest.UserController 에서 검사한다.
func LoginRequiredForWrites() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}

		p := strings.TrimSuffix(c.Request.URL.Path, "/")
		if !loginRequiredWritePaths[p] || (c.Request.Method == http.MethodPost && p == "/api/user") {
			c.Next()
			return
		}

		if user, ok := sessions.Default(c).Get("user").(*models.User); ok && user != nil && user.Id > 0 {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": "error", "message": "로그인이 필요합니다."})
	}
}

type loginRequest struct {
	Loginid string `json:"loginid" form:"loginid"`
	Passwd  string `json:"passwd" form:"passwd"`
}

// SetAuthRouter 는 로그인/로그아웃/내 정보 라우트를 등록한다 (수동 관리 — buildtool-router 재생성 대상 아님).
// 로그인은 POST 바디(JSON 또는 form)로만 받는다. 쿼리스트링에 비밀번호를 싣지 않기 위함.
func SetAuthRouter(r *gin.Engine) {
	r.POST("/api/login/login", func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "error", "message": "invalid request"})
			return
		}
		var controller api.Login
		controller.Init(c)
		controller.Login(req.Loginid, req.Passwd)
		controller.Close()
		c.JSON(controller.Code, controller.Result)
	})

	r.POST("/api/login/logout", func(c *gin.Context) {
		var controller api.Login
		controller.Init(c)
		controller.Logout()
		controller.Close()
		c.JSON(controller.Code, controller.Result)
	})

	r.GET("/api/login/me", func(c *gin.Context) {
		var controller api.Login
		controller.Init(c)
		controller.Me()
		controller.Close()
		c.JSON(controller.Code, controller.Result)
	})
}
