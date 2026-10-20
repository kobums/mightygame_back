package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

// PlayerTokenHeader 는 게임 진행 API 에서 플레이어 신원을 증명하는 헤더 이름이다.
// 방 생성/입장 시 서버가 발급한 불투명 토큰을 이 헤더로 보내야 하며,
// 서버는 클라이언트가 보낸 playerId 를 믿지 않고 토큰에서 playerId 를 도출한다.
const PlayerTokenHeader = "X-Player-Token"

// playerBinding 은 토큰이 묶인 방과 플레이어.
type playerBinding struct {
	RoomID   int64
	PlayerID int64
}

var (
	playerTokens      = make(map[string]playerBinding)
	playerTokensMutex sync.RWMutex
)

// issuePlayerToken 은 (roomID, playerID) 에 묶인 랜덤 토큰을 발급해 서버 메모리에 저장한다.
func issuePlayerToken(roomID int64, playerID int64) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	playerTokensMutex.Lock()
	playerTokens[token] = playerBinding{RoomID: roomID, PlayerID: playerID}
	playerTokensMutex.Unlock()
	return token, nil
}

// revokePlayerToken 은 토큰 하나를 폐기한다.
func revokePlayerToken(token string) {
	playerTokensMutex.Lock()
	delete(playerTokens, token)
	playerTokensMutex.Unlock()
}

// revokeRoomTokens 는 방에 묶인 모든 토큰을 폐기한다 (방 삭제 시).
func revokeRoomTokens(roomID int64) {
	playerTokensMutex.Lock()
	for token, binding := range playerTokens {
		if binding.RoomID == roomID {
			delete(playerTokens, token)
		}
	}
	playerTokensMutex.Unlock()
}

// playerToken 은 요청에서 플레이어 토큰을 꺼낸다.
// 헤더를 우선하고, 헤더를 실을 수 없는 navigator.sendBeacon(방 나가기) 용으로 본문 값을 허용한다.
func (c *RoomController) playerToken(bodyToken string) string {
	if token := c.Context.GetHeader(PlayerTokenHeader); token != "" {
		return token
	}
	return bodyToken
}

// authPlayer 는 토큰을 검증해 해당 방에서의 playerId 를 돌려준다.
// 토큰이 없거나, 모르는 토큰이거나, 다른 방의 토큰이면 401 을 설정하고 false 를 돌려준다.
func (c *RoomController) authPlayer(roomID int64, token string) (int64, bool) {
	if token != "" {
		playerTokensMutex.RLock()
		binding, ok := playerTokens[token]
		playerTokensMutex.RUnlock()
		if ok && binding.RoomID == roomID {
			return binding.PlayerID, true
		}
	}

	c.Code = http.StatusUnauthorized
	c.Result["code"] = "error"
	c.Result["message"] = "플레이어 인증에 실패했습니다. 방에 다시 입장해 주세요"
	return 0, false
}
