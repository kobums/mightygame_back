package api

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"mighty/controllers"
	"mighty/models"

	"github.com/gin-contrib/sessions"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	controllers.Controller
}

// 존재하지 않는 아이디일 때도 bcrypt 비교를 한 번 수행하기 위한 더미 해시.
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("mighty-dummy-password"), bcrypt.DefaultCost)

// HashPasswd 는 비밀번호를 bcrypt 로 해시한다.
func HashPasswd(passwd string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// IsHashed 는 저장된 값이 bcrypt 해시인지 판단한다.
func IsHashed(stored string) bool {
	return strings.HasPrefix(stored, "$2a$") || strings.HasPrefix(stored, "$2b$") || strings.HasPrefix(stored, "$2y$")
}

// verifyPasswd 는 저장된 값이 bcrypt 해시면 bcrypt 로, 아니면(기존 평문 계정) 상수 시간 비교로 검증한다.
// 평문 계정이 일치하면 true 와 함께 migrate=true 를 돌려줘 호출 측에서 bcrypt 로 바꿔 저장하게 한다.
func verifyPasswd(stored string, raw string) (ok bool, migrate bool) {
	if IsHashed(stored) {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(raw)) == nil, false
	}
	if stored == "" {
		return false, false
	}
	if subtle.ConstantTimeCompare([]byte(stored), []byte(raw)) == 1 {
		return true, true
	}
	return false, false
}

func (c *Login) fail() {
	c.Code = http.StatusUnauthorized
	c.Set("code", "error")
	// 아이디 존재 여부를 알 수 없도록 같은 메시지를 쓴다.
	c.Set("message", "아이디 또는 비밀번호가 올바르지 않습니다.")
}

// Login 은 아이디/비밀번호를 검증하고 세션 쿠키에 사용자 정보를 저장한다.
func (c *Login) Login(loginid string, passwd string) {
	if loginid == "" || passwd == "" {
		c.fail()
		return
	}

	manager := models.NewUserManager(c.NewConnection())
	var args []interface{}
	args = append(args, models.Where{Column: "loginid", Value: loginid, Compare: "="})
	items := manager.Find(args)
	if items == nil || len(*items) != 1 {
		// 비교 비용을 맞춰 응답 시간으로 계정 존재 여부를 추측하기 어렵게 한다.
		bcrypt.CompareHashAndPassword(dummyHash, []byte(passwd))
		c.fail()
		return
	}
	user := (*items)[0]

	ok, migrate := verifyPasswd(user.Passwd, passwd)
	if !ok {
		c.fail()
		return
	}

	if migrate {
		// 평문으로 저장돼 있던 계정은 로그인 성공 시 bcrypt 로 바꿔 저장한다.
		if hashed, err := HashPasswd(passwd); err == nil {
			user.Passwd = hashed
			if err := manager.Update(&user); err != nil {
				log.Println("password migration failed:", err)
			}
		}
	}

	// 쿠키 세션은 서명만 되고 암호화되지 않으므로 비밀번호(해시)는 넣지 않는다.
	sessionUser := &models.User{Id: user.Id, Loginid: user.Loginid, Name: user.Name, Date: user.Date}
	session := sessions.Default(c.Context)
	session.Set("user", sessionUser)
	if err := session.Save(); err != nil {
		log.Println("session save failed:", err)
		c.Code = http.StatusInternalServerError
		c.Set("code", "error")
		return
	}

	c.Session = sessionUser
	c.Set("session", sessionUser)
	c.Set("user", sessionUser)
}

// Logout 은 세션을 비운다.
func (c *Login) Logout() {
	session := sessions.Default(c.Context)
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1, HttpOnly: true})
	if err := session.Save(); err != nil {
		log.Println("session save failed:", err)
	}
	c.Session = nil
	c.Set("session", nil)
}

// Me 는 현재 로그인한 사용자를 돌려준다.
func (c *Login) Me() {
	if c.Session == nil {
		c.Code = http.StatusUnauthorized
		c.Set("code", "error")
		c.Set("message", "로그인이 필요합니다.")
		return
	}
	c.Set("user", c.Session)
}
