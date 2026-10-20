// 주의: 이 파일은 buildtool-router 로 생성된 뒤 수동으로 보안 수정이 들어간 파일이다.
// buildtool-router 의 소스/템플릿은 어느 저장소에도 없으므로(gomachine 의 views/go/router.jet 은
// fiber 용 router/routers/* 구조라 이 gin 파일과 다르다) 재생성하면 아래 수정이 모두 사라진다.
// 재생성했다면 다음을 다시 적용할 것:
//   - PUT 핸들러: GetRawData 로 읽은 본문을 json.Unmarshal 로 다시 디코드 (ShouldBindJSON 은 이미 소비된 본문을 못 읽음)
//   - GET /api/login/login?loginid=&passwd= 라우트를 만들지 말 것 (로그인은 router/auth.go 의 POST)
// 로그인 필수 미들웨어(router/auth.go LoginRequiredForWrites)·본인 확인(controllers/rest/user.go)·
// bcrypt 해시(controllers/api/login.go)·게임 플레이어 토큰(controllers/api/playertoken.go)은
// 이 파일 밖에 있어 재생성의 영향을 받지 않는다. 단 controllers/rest/*.go 도 재생성 대상이면 함께 확인할 것.

package router

import (
	"encoding/json"
    "strconv"
	"mighty/controllers/api"
	"mighty/controllers/rest"
	"mighty/models"
	"github.com/gin-gonic/gin"
)

func SetRouter(r *gin.Engine) {
	apiGroup := r.Group("/api")
	{

		// 로그인/로그아웃은 router/auth.go 의 SetAuthRouter (POST /api/login/login) 로 옮겼다.
		// GET /api/login/login?loginid=&passwd= 는 비밀번호가 URL·로그에 남으므로 제공하지 않는다
		// (기존 Login 은 빈 함수였고 이를 호출하는 클라이언트도 없음).

		// Room endpoints
		apiGroup.POST("/room/create", func(c *gin.Context) {
			var controller api.RoomController
			controller.Init(c)
			controller.CreateRoom()
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.GET("/room/list", func(c *gin.Context) {
			var controller api.RoomController
			controller.Init(c)
			controller.GetRoomList()
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.POST("/room/join", func(c *gin.Context) {
			var controller api.RoomController
			controller.Init(c)
			controller.JoinRoom()
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.POST("/room/leave", func(c *gin.Context) {
			var controller api.RoomController
			controller.Init(c)
			controller.LeaveRoom()
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.GET("/room/:id", func(c *gin.Context) {
			var controller api.RoomController
			controller.Init(c)
			controller.GetRoomDetail()
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.POST("/room/start", func(c *gin.Context) {
			var controller api.RoomController
			controller.Init(c)
			controller.StartGame()
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		// Game state endpoints
		apiGroup.GET("/game/state/:roomId", func(c *gin.Context) {
			var controller api.RoomController
			controller.Init(c)
			controller.GetGameState()
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.POST("/game/bid", func(c *gin.Context) {
			var controller api.RoomController
			controller.Init(c)
			controller.PlaceBid()
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

	}

	{

		apiGroup.GET("/game/:id", func(c *gin.Context) {
			id_, _ := strconv.ParseInt(c.Param("id"), 10, 64)
			var controller rest.GameController
			controller.Init(c)
			controller.Read(id_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.GET("/game", func(c *gin.Context) {
			page_, _ := strconv.Atoi(c.Query("page"))
			pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
			var controller rest.GameController
			controller.Init(c)
			controller.Index(page_, pagesize_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.POST("/game", func(c *gin.Context) {
			var item_ models.Game
			c.ShouldBindJSON(&item_)
			var controller rest.GameController
			controller.Init(c)
			controller.Insert(item_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.PUT("/game", func(c *gin.Context) {
			var results map[string]interface{}
			jsonData, _ := c.GetRawData()
			json.Unmarshal(jsonData, &results)
			var item_ models.Game
			// GetRawData 로 본문을 이미 읽었으므로 ShouldBindJSON 대신 같은 데이터로 디코드한다.
			json.Unmarshal(jsonData, &item_)
			var controller rest.GameController
			controller.Init(c)
			controller.Update(item_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.DELETE("/game", func(c *gin.Context) {
			var item_ models.Game
			c.ShouldBindJSON(&item_)
			var controller rest.GameController
			controller.Init(c)
			controller.Delete(item_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.GET("/gameuser/:id", func(c *gin.Context) {
			id_, _ := strconv.ParseInt(c.Param("id"), 10, 64)
			var controller rest.GameuserController
			controller.Init(c)
			controller.Read(id_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.GET("/gameuser", func(c *gin.Context) {
			page_, _ := strconv.Atoi(c.Query("page"))
			pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
			var controller rest.GameuserController
			controller.Init(c)
			controller.Index(page_, pagesize_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.POST("/gameuser", func(c *gin.Context) {
			var item_ models.Gameuser
			c.ShouldBindJSON(&item_)
			var controller rest.GameuserController
			controller.Init(c)
			controller.Insert(item_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.PUT("/gameuser", func(c *gin.Context) {
			var results map[string]interface{}
			jsonData, _ := c.GetRawData()
			json.Unmarshal(jsonData, &results)
			var item_ models.Gameuser
			// GetRawData 로 본문을 이미 읽었으므로 ShouldBindJSON 대신 같은 데이터로 디코드한다.
			json.Unmarshal(jsonData, &item_)
			var controller rest.GameuserController
			controller.Init(c)
			controller.Update(item_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.DELETE("/gameuser", func(c *gin.Context) {
			var item_ models.Gameuser
			c.ShouldBindJSON(&item_)
			var controller rest.GameuserController
			controller.Init(c)
			controller.Delete(item_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.GET("/user/:id", func(c *gin.Context) {
			id_, _ := strconv.ParseInt(c.Param("id"), 10, 64)
			var controller rest.UserController
			controller.Init(c)
			controller.Read(id_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.GET("/user", func(c *gin.Context) {
			page_, _ := strconv.Atoi(c.Query("page"))
			pagesize_, _ := strconv.Atoi(c.Query("pagesize"))
			var controller rest.UserController
			controller.Init(c)
			controller.Index(page_, pagesize_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.POST("/user", func(c *gin.Context) {
			var item_ models.User
			c.ShouldBindJSON(&item_)
			var controller rest.UserController
			controller.Init(c)
			controller.Insert(item_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.PUT("/user", func(c *gin.Context) {
			var results map[string]interface{}
			jsonData, _ := c.GetRawData()
			json.Unmarshal(jsonData, &results)
			var item_ models.User
			// GetRawData 로 본문을 이미 읽었으므로 ShouldBindJSON 대신 같은 데이터로 디코드한다.
			json.Unmarshal(jsonData, &item_)
			var controller rest.UserController
			controller.Init(c)
			controller.Update(item_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

		apiGroup.DELETE("/user", func(c *gin.Context) {
			var item_ models.User
			c.ShouldBindJSON(&item_)
			var controller rest.UserController
			controller.Init(c)
			controller.Delete(item_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

	}

}
