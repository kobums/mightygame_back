package rest

import (
	"regexp"
	"mighty/controllers"
	"mighty/models"
)

type GameuserController struct {
	controllers.Controller
}

func (c *GameuserController) Read(id int64) {
	conn := c.NewConnection()

	manager := models.NewGameuserManager(conn)
	item := manager.Get(id)
	c.Set("item", item)
}

func (c *GameuserController) Index(page int, pagesize int) {
	conn := c.NewConnection()

	manager := models.NewGameuserManager(conn)

    var args []interface{}
    
    _game := c.Geti64("game")
    if _game != 0 {
        args = append(args, models.Where{Column:"game", Value:_game, Compare:"="})    
    }
    _user := c.Geti64("user")
    if _user != 0 {
        args = append(args, models.Where{Column:"user", Value:_user, Compare:"="})    
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
	c.Set("items", items)

    total := manager.Count(args)
	c.Set("total", total)
}

func (c *GameuserController) Insert(item models.Gameuser) {
	conn := c.NewConnection()

	manager := models.NewGameuserManager(conn)
	manager.Insert(&item)

    c.Result["id"] = manager.GetIdentity()
}

func (c *GameuserController) Update(item models.Gameuser) {
	conn := c.NewConnection()

	manager := models.NewGameuserManager(conn)
	manager.Update(&item)
}

func (c *GameuserController) Delete(item models.Gameuser) {
	conn := c.NewConnection()

	manager := models.NewGameuserManager(conn)
	manager.Delete(item.Id)
}
