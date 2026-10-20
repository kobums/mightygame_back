package rest

import (
	"net/http"
	"regexp"

	"mighty/controllers"
	"mighty/controllers/api"
	"mighty/models"
)

type UserController struct {
	controllers.Controller
}

func (c *UserController) Read(id int64) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	item := manager.Get(id)
	// 응답에 비밀번호를 노출하지 않는다.
	if item != nil {
		item.Passwd = ""
	}
	c.Set("item", item)
}

func (c *UserController) Index(page int, pagesize int) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)

    var args []interface{}
    
    _loginid := c.Get("loginid")
    if _loginid != "" {
        args = append(args, models.Where{Column:"loginid", Value:_loginid, Compare:"like"})
    }
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
    }
    _startdate := c.Get("startdate")
    _enddate := c.Get("enddate")
    if _startdate != "" && _enddate != "" {        
        var v [2]string
        v[0] = _startdate
        v[1] = _enddate  
        args = append(args, models.Where{Column:"date", Value:v, Compare:"between"})    
    } else if  _startdate != "" {          
        args = append(args, models.Where{Column:"date", Value:_startdate, Compare:">="})
    } else if  _enddate != "" {          
        args = append(args, models.Where{Column:"date", Value:_enddate, Compare:"<="})            
    }
    

    if page != 0 && pagesize != 0 {
        args = append(args, models.Paging(page, pagesize))
    }
    
    orderby := c.Get("orderby")
    if orderby == "" {
        if page != 0 && pagesize != 0 {
            orderby = "id desc"
        }
    }

    // orderby 는 ORDER BY 절에 그대로 들어가므로 "컬럼 [asc|desc]" 목록 형태만 허용한다 (SQL injection 방지).
    if orderby != "" && regexp.MustCompile(`^\s*[A-Za-z0-9_]+(\s+(?i:asc|desc))?\s*(,\s*[A-Za-z0-9_]+(\s+(?i:asc|desc))?\s*)*$`).MatchString(orderby) {
        args = append(args, models.Ordering(orderby))
    }
    
	items := manager.Find(args)
	// 응답에 비밀번호를 노출하지 않는다.
	if items != nil {
		for i := range *items {
			(*items)[i].Passwd = ""
		}
	}
	c.Set("items", items)

    total := manager.Count(args)
	c.Set("total", total)
}

func (c *UserController) forbidden() {
	c.Code = http.StatusForbidden
	c.Set("code", "error")
	c.Set("message", "본인 계정만 수정/삭제할 수 있습니다.")
}

func (c *UserController) badRequest(message string) {
	c.Code = http.StatusBadRequest
	c.Set("code", "error")
	c.Set("message", message)
}

// loginidTaken 은 다른 사용자(exceptId 제외)가 이미 쓰는 아이디인지 확인한다.
func loginidTaken(manager *models.UserManager, loginid string, exceptId int64) bool {
	var args []interface{}
	args = append(args, models.Where{Column: "loginid", Value: loginid, Compare: "="})
	items := manager.Find(args)
	if items == nil {
		return false
	}
	for _, v := range *items {
		if v.Id != exceptId {
			return true
		}
	}
	return false
}

// Insert 는 회원가입. 로그인 없이 호출할 수 있으며 비밀번호는 bcrypt 로 해시해 저장한다.
func (c *UserController) Insert(item models.User) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)

	if item.Loginid == "" || item.Passwd == "" {
		c.badRequest("아이디와 비밀번호를 입력하세요.")
		return
	}
	if loginidTaken(manager, item.Loginid, 0) {
		c.badRequest("이미 사용 중인 아이디입니다.")
		return
	}

	hashed, err := api.HashPasswd(item.Passwd)
	if err != nil {
		c.Code = http.StatusInternalServerError
		c.Set("code", "error")
		return
	}
	item.Passwd = hashed
	// id 는 서버가 정한다 (클라이언트가 임의 id 로 가입하지 못하게).
	item.Id = 0
	item.Date = ""

	if err := manager.Insert(&item); err != nil {
		c.Code = http.StatusInternalServerError
		c.Set("code", "error")
		return
	}

	c.Result["id"] = manager.GetIdentity()
}

// Update 는 로그인한 사용자 본인의 계정만 수정한다.
// 비밀번호를 보내지 않으면 기존 비밀번호를 유지하고, 보내면 bcrypt 로 해시한다.
func (c *UserController) Update(item models.User) {
	if c.Session == nil || item.Id == 0 || item.Id != c.Session.Id {
		c.forbidden()
		return
	}

	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	existing := manager.Get(item.Id)
	if existing == nil {
		c.forbidden()
		return
	}

	if item.Loginid == "" {
		item.Loginid = existing.Loginid
	} else if item.Loginid != existing.Loginid && loginidTaken(manager, item.Loginid, item.Id) {
		c.badRequest("이미 사용 중인 아이디입니다.")
		return
	}
	if item.Name == "" {
		item.Name = existing.Name
	}
	if item.Date == "" {
		item.Date = existing.Date
	}
	if item.Passwd == "" {
		item.Passwd = existing.Passwd
	} else {
		hashed, err := api.HashPasswd(item.Passwd)
		if err != nil {
			c.Code = http.StatusInternalServerError
			c.Set("code", "error")
			return
		}
		item.Passwd = hashed
	}

	if err := manager.Update(&item); err != nil {
		c.Code = http.StatusInternalServerError
		c.Set("code", "error")
	}
}

// Delete 는 로그인한 사용자 본인의 계정만 삭제한다.
func (c *UserController) Delete(item models.User) {
	if c.Session == nil || item.Id == 0 || item.Id != c.Session.Id {
		c.forbidden()
		return
	}

	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	if err := manager.Delete(item.Id); err != nil {
		c.Code = http.StatusInternalServerError
		c.Set("code", "error")
	}
}
