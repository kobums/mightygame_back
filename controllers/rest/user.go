package rest

import (
	"regexp"
	"mighty/controllers"
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

func (c *UserController) Insert(item models.User) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	manager.Insert(&item)

    c.Result["id"] = manager.GetIdentity()
}

func (c *UserController) Update(item models.User) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	manager.Update(&item)
}

func (c *UserController) Delete(item models.User) {
	conn := c.NewConnection()

	manager := models.NewUserManager(conn)
	manager.Delete(item.Id)
}
