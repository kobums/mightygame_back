package rest

import (
	"mighty/controllers"
	"mighty/models"
)

type GameController struct {
	controllers.Controller
}

func (c *GameController) Read(id int64) {
	conn := c.NewConnection()

	manager := models.NewGameManager(conn)
	item := manager.Get(id)
	c.Set("item", item)
}

func (c *GameController) Index(page int, pagesize int) {
	conn := c.NewConnection()

	manager := models.NewGameManager(conn)

    var args []interface{}
    
    _name := c.Get("name")
    if _name != "" {
        args = append(args, models.Where{Column:"name", Value:_name, Compare:"="})
        
    }
    _member := c.Geti("member")
    if _member != 0 {
        args = append(args, models.Where{Column:"member", Value:_member, Compare:"="})    
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

    if orderby != "" {
        args = append(args, models.Ordering(orderby))
    }
    
	items := manager.Find(args)
	c.Set("items", items)

    total := manager.Count(args)
	c.Set("total", total)
}

func (c *GameController) Insert(item models.Game) {
	conn := c.NewConnection()

	manager := models.NewGameManager(conn)
	manager.Insert(&item)

    c.Result["id"] = manager.GetIdentity()
}

func (c *GameController) Update(item models.Game) {
	conn := c.NewConnection()

	manager := models.NewGameManager(conn)
	manager.Update(&item)
}

func (c *GameController) Delete(item models.Game) {
	conn := c.NewConnection()

	manager := models.NewGameManager(conn)
	manager.Delete(item.Id)
}
