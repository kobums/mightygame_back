package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

func newRoomTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handle := func(fn func(*RoomController)) gin.HandlerFunc {
		return func(g *gin.Context) {
			var c RoomController
			c.Init(g)
			fn(&c)
			c.Close()
			g.JSON(c.Code, c.Result)
		}
	}
	r.POST("/api/room/create", handle((*RoomController).CreateRoom))
	r.POST("/api/room/join", handle((*RoomController).JoinRoom))
	r.POST("/api/room/leave", handle((*RoomController).LeaveRoom))
	r.POST("/api/room/start", handle((*RoomController).StartGame))
	r.GET("/api/game/state/:roomId", handle((*RoomController).GetGameState))
	r.POST("/api/game/bid", handle((*RoomController).PlaceBid))
	return r
}

func doJSON(t *testing.T, r *gin.Engine, method, path, token string, body interface{}) (int, map[string]interface{}) {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set(PlayerTokenHeader, token)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	var out map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &out)
	return w.Code, out
}

func TestPlayerTokenBindsIdentity(t *testing.T) {
	r := newRoomTestRouter()

	_, created := doJSON(t, r, "POST", "/api/room/create", "", map[string]interface{}{"roomName": "r", "playerName": "host"})
	roomID := int64(created["roomId"].(float64))
	hostToken, _ := created["playerToken"].(string)
	if hostToken == "" {
		t.Fatalf("no token on create: %v", created)
	}

	tokens := []string{hostToken}
	for i := 0; i < 4; i++ {
		_, joined := doJSON(t, r, "POST", "/api/room/join", "", map[string]interface{}{"roomId": roomID, "playerName": "p" + strconv.Itoa(i)})
		tok, _ := joined["playerToken"].(string)
		if tok == "" {
			t.Fatalf("no token on join: %v", joined)
		}
		tokens = append(tokens, tok)
	}

	// 방장이 아닌 플레이어는 시작할 수 없다 (본문 playerId 로 방장 사칭 불가).
	if code, _ := doJSON(t, r, "POST", "/api/room/start", tokens[1], map[string]interface{}{"roomId": roomID, "playerId": created["playerId"]}); code != http.StatusForbidden {
		t.Fatalf("non-host start: got %d", code)
	}
	// 토큰 없이 시작 불가.
	if code, _ := doJSON(t, r, "POST", "/api/room/start", "", map[string]interface{}{"roomId": roomID}); code != http.StatusUnauthorized {
		t.Fatalf("tokenless start: got %d", code)
	}
	if code, out := doJSON(t, r, "POST", "/api/room/start", hostToken, map[string]interface{}{"roomId": roomID}); code != http.StatusOK || out["code"] != "success" {
		t.Fatalf("host start: %d %v", code, out)
	}

	statePath := "/api/game/state/" + strconv.FormatInt(roomID, 10)
	// 토큰 없이 ?playerId= 로 남의 패를 볼 수 없다.
	if code, _ := doJSON(t, r, "GET", statePath+"?playerId=1", "", nil); code != http.StatusUnauthorized {
		t.Fatalf("tokenless state: got %d", code)
	}
	code, state := doJSON(t, r, "GET", statePath, tokens[2], nil)
	if code != http.StatusOK {
		t.Fatalf("state: %d %v", code, state)
	}
	game := state["game"].(map[string]interface{})
	if hand, _ := game["myHand"].([]interface{}); len(hand) != 10 {
		t.Fatalf("expected 10 cards, got %v", game["myHand"])
	}

	// 차례가 아닌 플레이어가 본문 playerId 로 현재 비더를 사칭해도 거부된다.
	_, out := doJSON(t, r, "POST", "/api/game/bid", tokens[3], map[string]interface{}{"roomId": roomID, "playerId": game["currentBidder"], "bidType": "13", "trumpSuit": "spade"})
	if out["code"] != "error" {
		t.Fatalf("impersonated bid accepted: %v", out)
	}
	_, out = doJSON(t, r, "POST", "/api/game/bid", hostToken, map[string]interface{}{"roomId": roomID, "bidType": "13", "trumpSuit": "spade"})
	if out["code"] != "success" {
		t.Fatalf("host bid: %v", out)
	}

	// 다른 방의 토큰은 쓸 수 없다.
	if code, _ := doJSON(t, r, "GET", "/api/game/state/999999", hostToken, nil); code != http.StatusUnauthorized {
		t.Fatalf("cross-room token: got %d", code)
	}

	// 게임 중 자동 나가기는 게임을 깨지 않는다.
	doJSON(t, r, "POST", "/api/room/leave", "", map[string]interface{}{"roomId": roomID, "playerToken": hostToken})
	if code, _ := doJSON(t, r, "GET", statePath, hostToken, nil); code != http.StatusOK {
		t.Fatalf("state after in-game leave: got %d", code)
	}
}

func TestLeaveRevokesToken(t *testing.T) {
	r := newRoomTestRouter()
	_, created := doJSON(t, r, "POST", "/api/room/create", "", map[string]interface{}{"roomName": "r", "playerName": "host"})
	roomID := int64(created["roomId"].(float64))
	_, joined := doJSON(t, r, "POST", "/api/room/join", "", map[string]interface{}{"roomId": roomID, "playerName": "guest"})
	guestToken := joined["playerToken"].(string)

	// 게스트가 본문 playerId 로 방장을 내보낼(=방 삭제) 수 없다: 자기 자신만 나간다.
	_, out := doJSON(t, r, "POST", "/api/room/leave", guestToken, map[string]interface{}{"roomId": roomID, "playerId": created["playerId"]})
	if out["roomDeleted"] == true {
		t.Fatalf("guest deleted room via host playerId: %v", out)
	}
	if code, _ := doJSON(t, r, "POST", "/api/room/leave", guestToken, map[string]interface{}{"roomId": roomID}); code != http.StatusUnauthorized {
		t.Fatalf("revoked token reused: got %d", code)
	}
}
