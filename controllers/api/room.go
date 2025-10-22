package api

import (
	"mighty/controllers"
	"strconv"
	"sync"
	"time"
)

type RoomController struct {
	controllers.Controller
}

// Player represents a player in a room
type Player struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Room represents a game room
type Room struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	HostID         int64     `json:"hostId"`
	HostName       string    `json:"hostName"`
	Players        []Player  `json:"players"`
	CurrentPlayers int       `json:"currentPlayers"`
	MaxPlayers     int       `json:"maxPlayers"`
	Status         string    `json:"status"` // "waiting", "playing"
	CreatedAt      time.Time `json:"createdAt"`
}

// In-memory storage
var (
	rooms      = make(map[int64]*Room)
	roomsMutex sync.RWMutex
	nextRoomID int64 = 1
	nextPlayerID int64 = 1
)

// CreateRoom creates a new game room
func (c *RoomController) CreateRoom() {
	// Parse JSON body
	var reqBody struct {
		RoomName   string `json:"roomName"`
		PlayerName string `json:"playerName"`
	}
	c.Context.ShouldBindJSON(&reqBody)

	roomName := reqBody.RoomName
	playerName := reqBody.PlayerName

	if roomName == "" || playerName == "" {
		c.Result["code"] = "error"
		c.Result["message"] = "방 이름과 플레이어 이름은 필수입니다"
		return
	}

	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	// Create new room
	roomID := nextRoomID
	nextRoomID++

	playerID := nextPlayerID
	nextPlayerID++

	room := &Room{
		ID:             roomID,
		Name:           roomName,
		HostID:         playerID,
		HostName:       playerName,
		Players:        []Player{{ID: playerID, Name: playerName}},
		CurrentPlayers: 1,
		MaxPlayers:     5,
		Status:         "waiting",
		CreatedAt:      time.Now(),
	}

	rooms[roomID] = room

	c.Result["code"] = "success"
	c.Result["roomId"] = roomID
	c.Result["playerId"] = playerID
	c.Result["room"] = room
}

// GetRoomList returns list of available rooms
func (c *RoomController) GetRoomList() {
	roomsMutex.RLock()
	defer roomsMutex.RUnlock()

	var roomList []Room
	for _, room := range rooms {
		if room.Status == "waiting" && room.CurrentPlayers < room.MaxPlayers {
			roomList = append(roomList, *room)
		}
	}

	c.Result["code"] = "success"
	c.Result["rooms"] = roomList
}

// JoinRoom adds a player to an existing room
func (c *RoomController) JoinRoom() {
	var reqBody struct {
		RoomID     int64  `json:"roomId"`
		PlayerName string `json:"playerName"`
	}
	c.Context.ShouldBindJSON(&reqBody)

	roomID := reqBody.RoomID
	playerName := reqBody.PlayerName

	if playerName == "" {
		c.Result["code"] = "error"
		c.Result["message"] = "플레이어 이름은 필수입니다"
		return
	}

	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		c.Result["code"] = "error"
		c.Result["message"] = "방을 찾을 수 없습니다"
		return
	}

	if room.CurrentPlayers >= room.MaxPlayers {
		c.Result["code"] = "error"
		c.Result["message"] = "방이 가득 찼습니다"
		return
	}

	if room.Status != "waiting" {
		c.Result["code"] = "error"
		c.Result["message"] = "게임이 이미 시작되었습니다"
		return
	}

	playerID := nextPlayerID
	nextPlayerID++

	room.Players = append(room.Players, Player{ID: playerID, Name: playerName})
	room.CurrentPlayers++

	c.Result["code"] = "success"
	c.Result["playerId"] = playerID
	c.Result["room"] = room
}

// LeaveRoom removes a player from a room
func (c *RoomController) LeaveRoom() {
	var reqBody struct {
		RoomID   int64 `json:"roomId"`
		PlayerID int64 `json:"playerId"`
	}
	c.Context.ShouldBindJSON(&reqBody)

	roomID := reqBody.RoomID
	playerID := reqBody.PlayerID

	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		c.Result["code"] = "error"
		c.Result["message"] = "방을 찾을 수 없습니다"
		return
	}

	// Find and remove player
	playerIndex := -1
	for i, player := range room.Players {
		if player.ID == playerID {
			playerIndex = i
			break
		}
	}

	if playerIndex == -1 {
		c.Result["code"] = "error"
		c.Result["message"] = "플레이어를 찾을 수 없습니다"
		return
	}

	// Remove player
	room.Players = append(room.Players[:playerIndex], room.Players[playerIndex+1:]...)
	room.CurrentPlayers--

	// If room is empty or host left, delete room
	if room.CurrentPlayers == 0 || playerID == room.HostID {
		delete(rooms, roomID)
		c.Result["code"] = "success"
		c.Result["message"] = "방을 나갔습니다"
		c.Result["roomDeleted"] = true
		return
	}

	c.Result["code"] = "success"
	c.Result["message"] = "방을 나갔습니다"
}

// GetRoomDetail returns detailed room information
func (c *RoomController) GetRoomDetail() {
	roomID, _ := strconv.ParseInt(c.Context.Param("id"), 10, 64)

	roomsMutex.RLock()
	defer roomsMutex.RUnlock()

	room, exists := rooms[roomID]
	if !exists {
		c.Result["code"] = "error"
		c.Result["message"] = "방을 찾을 수 없습니다"
		return
	}

	c.Result["code"] = "success"
	c.Result["room"] = room
}

// StartGame starts the game
func (c *RoomController) StartGame() {
	var reqBody struct {
		RoomID   int64 `json:"roomId"`
		PlayerID int64 `json:"playerId"`
	}
	c.Context.ShouldBindJSON(&reqBody)

	roomID := reqBody.RoomID
	_ = reqBody.PlayerID // TODO: Use this to verify player is host

	roomsMutex.Lock()
	defer roomsMutex.Unlock()

	room, exists := rooms[roomID]
	if !exists {
		c.Result["code"] = "error"
		c.Result["message"] = "방을 찾을 수 없습니다"
		return
	}

	if room.CurrentPlayers != 5 {
		c.Result["code"] = "error"
		c.Result["message"] = "5명의 플레이어가 필요합니다"
		return
	}

	room.Status = "playing"

	c.Result["code"] = "success"
	c.Result["message"] = "게임이 시작되었습니다"
	c.Result["room"] = room
}
