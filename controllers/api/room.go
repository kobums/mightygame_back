package api

import (
	"math/rand"
	"mighty/controllers"
	"strconv"
	"sync"
	"time"
)

type RoomController struct {
	controllers.Controller
}

// Card represents a playing card
type Card struct {
	Suit   string `json:"suit"`   // "spade", "diamond", "heart", "club", "joker"
	Number string `json:"number"` // "A", "2"..."10", "J", "Q", "K", "JOKER"
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

// GameState represents the current game state
type GameState struct {
	RoomID         int64              `json:"roomId"`
	Players        []Player           `json:"players"`
	PlayerHands    map[int64][]Card   `json:"-"` // Hidden from client
	Kitty          []Card             `json:"-"` // Hidden from client
	CurrentTurn    int                `json:"currentTurn"`    // Index of player whose turn it is
	Phase          string             `json:"phase"`          // "bidding", "kitty_exchange", "playing"
	Bids           []Bid              `json:"bids"`
	CurrentBidder  int64              `json:"currentBidder"`  // PlayerID of current bidder
	HighestBid     *Bid               `json:"highestBid"`
	TrumpSuit      string             `json:"trumpSuit"`
}

// Bid represents a bidding action
type Bid struct {
	PlayerID   int64  `json:"playerId"`
	PlayerName string `json:"playerName"`
	BidType    string `json:"bidType"` // "pass", "13", "14", "15", "16", "17", "18", "19", "20"
	TrumpSuit  string `json:"trumpSuit,omitempty"` // For actual bids (not pass)
}

// In-memory storage
var (
	rooms        = make(map[int64]*Room)
	games        = make(map[int64]*GameState)
	roomsMutex   sync.RWMutex
	gamesMutex   sync.RWMutex
	nextRoomID   int64 = 1
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

// createDeck creates a standard 53-card deck (52 cards + 1 joker)
func createDeck() []Card {
	suits := []string{"spade", "diamond", "heart", "club"}
	numbers := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	deck := make([]Card, 0, 53)

	// Add all regular cards
	for _, suit := range suits {
		for _, number := range numbers {
			deck = append(deck, Card{Suit: suit, Number: number})
		}
	}

	// Add joker
	deck = append(deck, Card{Suit: "joker", Number: "JOKER"})

	return deck
}

// shuffleDeck shuffles the deck using Fisher-Yates algorithm
func shuffleDeck(deck []Card) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(deck) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

// initializeGame creates initial game state with shuffled and dealt cards
func initializeGame(room *Room) *GameState {
	// Create and shuffle deck
	deck := createDeck()
	shuffleDeck(deck)

	// Deal 10 cards to each player
	playerHands := make(map[int64][]Card)
	cardIndex := 0
	for _, player := range room.Players {
		playerHands[player.ID] = deck[cardIndex : cardIndex+10]
		cardIndex += 10
	}

	// Remaining 3 cards go to kitty
	kitty := deck[50:53]

	game := &GameState{
		RoomID:        room.ID,
		Players:       room.Players,
		PlayerHands:   playerHands,
		Kitty:         kitty,
		CurrentTurn:   0,
		Phase:         "bidding",
		Bids:          []Bid{},
		CurrentBidder: room.Players[0].ID,
		HighestBid:    nil,
		TrumpSuit:     "",
	}

	return game
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
	room, exists := rooms[roomID]
	if !exists {
		roomsMutex.Unlock()
		c.Result["code"] = "error"
		c.Result["message"] = "방을 찾을 수 없습니다"
		return
	}

	if room.CurrentPlayers != 5 {
		roomsMutex.Unlock()
		c.Result["code"] = "error"
		c.Result["message"] = "5명의 플레이어가 필요합니다"
		return
	}

	room.Status = "playing"
	roomsMutex.Unlock()

	// Initialize game state
	gamesMutex.Lock()
	games[roomID] = initializeGame(room)
	gamesMutex.Unlock()

	c.Result["code"] = "success"
	c.Result["message"] = "게임이 시작되었습니다"
	c.Result["room"] = room
}

// GetGameState returns current game state for a player
func (c *RoomController) GetGameState() {
	roomID, _ := strconv.ParseInt(c.Context.Param("roomId"), 10, 64)
	playerID, _ := strconv.ParseInt(c.Context.Query("playerId"), 10, 64)

	gamesMutex.RLock()
	game, exists := games[roomID]
	gamesMutex.RUnlock()

	if !exists {
		c.Result["code"] = "error"
		c.Result["message"] = "게임을 찾을 수 없습니다"
		return
	}

	// Return game state with player's cards
	c.Result["code"] = "success"
	c.Result["game"] = map[string]interface{}{
		"roomId":        game.RoomID,
		"players":       game.Players,
		"myHand":        game.PlayerHands[playerID],
		"currentTurn":   game.CurrentTurn,
		"phase":         game.Phase,
		"bids":          game.Bids,
		"currentBidder": game.CurrentBidder,
		"highestBid":    game.HighestBid,
		"trumpSuit":     game.TrumpSuit,
	}
}

// PlaceBid handles player bidding
func (c *RoomController) PlaceBid() {
	var reqBody struct {
		RoomID    int64  `json:"roomId"`
		PlayerID  int64  `json:"playerId"`
		BidType   string `json:"bidType"`   // "pass", "13", "14", etc.
		TrumpSuit string `json:"trumpSuit"` // For actual bids (not pass)
	}
	c.Context.ShouldBindJSON(&reqBody)

	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	game, exists := games[reqBody.RoomID]
	if !exists {
		c.Result["code"] = "error"
		c.Result["message"] = "게임을 찾을 수 없습니다"
		return
	}

	// Validate it's the player's turn
	if game.CurrentBidder != reqBody.PlayerID {
		c.Result["code"] = "error"
		c.Result["message"] = "당신의 차례가 아닙니다"
		return
	}

	// Validate bid is higher than current highest
	if reqBody.BidType != "pass" {
		bidValue, _ := strconv.Atoi(reqBody.BidType)
		if game.HighestBid != nil {
			highestValue, _ := strconv.Atoi(game.HighestBid.BidType)
			if bidValue <= highestValue {
				c.Result["code"] = "error"
				c.Result["message"] = "더 높은 공약을 선택해야 합니다"
				return
			}
		}
	}

	// Find player name
	var playerName string
	for _, player := range game.Players {
		if player.ID == reqBody.PlayerID {
			playerName = player.Name
			break
		}
	}

	// Record bid
	bid := Bid{
		PlayerID:   reqBody.PlayerID,
		PlayerName: playerName,
		BidType:    reqBody.BidType,
		TrumpSuit:  reqBody.TrumpSuit,
	}
	game.Bids = append(game.Bids, bid)

	// Update highest bid if not a pass
	if reqBody.BidType != "pass" {
		game.HighestBid = &bid
		game.TrumpSuit = reqBody.TrumpSuit
	}

	// Move to next player
	game.CurrentTurn = (game.CurrentTurn + 1) % len(game.Players)
	game.CurrentBidder = game.Players[game.CurrentTurn].ID

	// Check if bidding is complete (4 passes after a bid, or all pass)
	passCount := 0
	for i := len(game.Bids) - 1; i >= 0 && i >= len(game.Bids)-4; i-- {
		if game.Bids[i].BidType == "pass" {
			passCount++
		}
	}

	if game.HighestBid != nil && passCount >= 4 {
		// Bidding complete - winner exchanges kitty
		game.Phase = "kitty_exchange"
	} else if len(game.Bids) >= 5 && game.HighestBid == nil {
		// All passed - redeal (for now just stay in bidding)
		game.Phase = "bidding"
	}

	c.Result["code"] = "success"
	c.Result["message"] = "비딩이 완료되었습니다"
	c.Result["game"] = game
}
