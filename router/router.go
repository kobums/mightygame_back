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

		apiGroup.GET("/login/login", func(c *gin.Context) {
			loginid_ := c.Query("loginid")
			passwd_ := c.Query("passwd")
			var controller api.Login
			controller.Init(c)
			controller.Login(loginid_, passwd_)
			controller.Close()
			c.JSON(controller.Code, controller.Result)
		})

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
			c.ShouldBindJSON(&item_)
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
			c.ShouldBindJSON(&item_)
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
			c.ShouldBindJSON(&item_)
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
