package models

import (
    "mighty/config"
    
    "database/sql"
    "errors"
    "fmt"
    "strings"
    "time"

    log "github.com/sirupsen/logrus"    
    _ "github.com/go-sql-driver/mysql"    
)

type Gameuser struct {
            
    Id                int64 `json:"id"`         
    Game                int64 `json:"game"`         
    User                int64 `json:"user"`         
    Date                string `json:"date"` 
    
    Extra                    interface{} `json:"extra"`
}

type GameuserManager struct {
    Conn    *sql.DB
    Result  *sql.Result
    Index   string
}

func NewGameuserManager(conn *sql.DB) *GameuserManager {
    var item GameuserManager

    if conn == nil {
        item.Conn = NewConnection()
    } else {
        item.Conn = conn
    }

    item.Index = ""

    return &item
}

func (p *GameuserManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *GameuserManager) SetIndex(index string) {
    p.Index = index
}

func (p *GameuserManager) GetQuery() string {
    ret := ""

    str := "select gu_id, gu_game, gu_user, gu_date, u_id, u_loginid, u_passwd, u_name, u_date from gameuser_tb, user_tb "

    if p.Index == "" {
        ret = str
    } else {
        ret = str + " use index(" + p.Index + ") "
    }

    ret += "where 1=1 "
    
    ret += "and gu_user = u_id "
    

    return ret;
}

func (p *GameuserManager) GetQuerySelect() string {
    ret := ""

    str := "select count(*) from gameuser_tb "

    if p.Index == "" {
        ret = str
    } else {
        ret = str + " use index(" + p.Index + ") "
    }

    return ret;
}

func (p *GameuserManager) Insert(item *Gameuser) error {
    if p.Conn == nil {
        return errors.New("Connection Error")
    }

    if item.Date == "" {
        t := time.Now()
        item.Date = fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
    }

    query := ""
    var res sql.Result
    var err error
    if item.Id > 0 {
        query = "insert into gameuser_tb (gu_id, gu_game, gu_user, gu_date) values (?, ?, ?, ?)"
        res, err = p.Conn.Exec(query , item.Id, item.Game, item.User, item.Date)
    } else {
        query = "insert into gameuser_tb (gu_game, gu_user, gu_date) values (?, ?, ?)"
        res, err = p.Conn.Exec(query , item.Game, item.User, item.Date)
    }
    
    if err == nil {
        p.Result = &res
    } else {
        log.Println(err)
        p.Result = nil
    }

    return err
}
func (p *GameuserManager) Delete(id int64) error {
    if p.Conn == nil {
        return errors.New("Connection Error")
    }

    query := "delete from gameuser_tb where gu_id = ?"
    _, err := p.Conn.Exec(query, id)

    return err
}
func (p *GameuserManager) Update(item *Gameuser) error {
    if p.Conn == nil {
        return errors.New("Connection Error")
    }

	query := "update gameuser_tb set gu_game = ?, gu_user = ?, gu_date = ? where gu_id = ?"
	_, err := p.Conn.Exec(query , item.Game, item.User, item.Date, item.Id)

    return err
}

func (p *GameuserManager) GetIdentity() int64 {
    if p.Result == nil {
        return 0
    }

    id, err := (*p.Result).LastInsertId()

    if err != nil {
        return 0
    } else {
        return id
    }
}

func (p *GameuserManager) ReadRow(rows *sql.Rows) *Gameuser {
    var item Gameuser
    var err error

    var _user User
    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Game, &item.User, &item.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Name, &_user.Date)
    } else {
        return nil
    }

    if err != nil {
        return nil
    } else {

        item.Extra = map[string]interface{}{

            "user":     _user,

        }
        return &item
    }
}

func (p *GameuserManager) ReadRows(rows *sql.Rows) *[]Gameuser {
    var items []Gameuser

    for rows.Next() {
        var item Gameuser
var _user User
    
    
        err := rows.Scan(&item.Id, &item.Game, &item.User, &item.Date, &_user.Id, &_user.Loginid, &_user.Passwd, &_user.Name, &_user.Date)

        if err != nil {
           log.Printf("ReadRows error : %v\n", err)
           break
        }

        item.Extra = map[string]interface{}{

            "user":     _user,

        }
        items = append(items, item)
    }


     return &items
}

func (p *GameuserManager) Get(id int64) *Gameuser {
    if p.Conn == nil {
        return nil
    }

    query := p.GetQuery() + " and gu_id = ?"

    
    query += " and gu_user = u_id"
    
    
    rows, err := p.Conn.Query(query, id)

    if err != nil {
        log.Printf("query error : %v, %v\n", err, query)
        return nil
    }

    defer rows.Close()

    return p.ReadRow(rows)
}

func (p *GameuserManager) Count(args []interface{}) int {
    if p.Conn == nil {
        return 0
    }

    var params []interface{}
    query := p.GetQuerySelect() + " where 1=1 "

    for _, arg := range args {
        switch v := arg.(type) {
        case Where:
            item := v

            if item.Compare == "in" {
                query += " and gu_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
            } else if item.Compare == "between" {
                query += " and gu_" + item.Column + " between ? and ?"

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                query += " and gu_" + item.Column + " " + item.Compare + " ?"
                if item.Compare == "like" {
                    params = append(params, "%" + item.Value.(string) + "%")
                } else {
                    params = append(params, item.Value)                
                }
            }
        }
    }
    
    rows, err := QueryArray(p.Conn, query, params)

    if err != nil {
        log.Printf("query error : %v, %v\n", err, query)
        return 0
    }

    defer rows.Close()

    if !rows.Next() {
        return 0
    }

    cnt := 0
    err = rows.Scan(&cnt)

    if err != nil {
        return 0
    } else {
        return cnt
    }
}

func (p *GameuserManager) Find(args []interface{}) *[]Gameuser {
    if p.Conn == nil {
        var items []Gameuser
        return &items
    }

    var params []interface{}
    query := p.GetQuery()

    page := 0
    pagesize := 0
    orderby := ""
    
    for _, arg := range args {
        switch v := arg.(type) {
        case PagingType:
            item := v
            page = item.Page
            pagesize = item.Pagesize
            break
        case OrderingType:
            item := v
            orderby = item.Order
            break
        case LimitType:
            item := v
            page = 1
            pagesize = item.Limit
            break
        case Option:
            item := v
            if item.Limit > 0 {
                page = 1
                pagesize = item.Limit
            } else {
                page = item.Page
                pagesize = item.Pagesize                
            }
            orderby = item.Order
            break
        case Where:
            item := v

            if item.Compare == "in" {
                query += " and gu_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
            } else if item.Compare == "between" {
                query += " and gu_" + item.Column + " between ? and ?"

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                query += " and gu_" + item.Column + " " + item.Compare + " ?"
                if item.Compare == "like" {
                    params = append(params, "%" + item.Value.(string) + "%")
                } else {
                    params = append(params, item.Value)                
                }
            }
        }
    }
    
    startpage := (page - 1) * pagesize
    
    if page > 0 && pagesize > 0 {
        if orderby == "" {
            orderby = "gu_id desc"
        } else {
            orderby = "gu_" + orderby
        }
        query += " order by " + orderby
        if config.Database == "mysql" {
            query += " limit ? offset ?"
            params = append(params, pagesize)
            params = append(params, startpage)
        } else if config.Database == "mssql" || config.Database == "sqlserver" {
            query += "OFFSET ? ROWS FETCH NEXT ? ROWS ONLY"
            params = append(params, startpage)
            params = append(params, pagesize)
        }
    } else {
        if orderby == "" {
            orderby = "gu_id"
        } else {
            orderby = "gu_" + orderby
        }
        query += " order by " + orderby
    }

    rows, err := QueryArray(p.Conn, query, params)

    if err != nil {
        log.Printf("query error : %v, %v\n", err, query)
        var items []Gameuser
        return &items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}


func (p *GameuserManager) FindByGame(game int64, args ...interface{}) *[]Gameuser {
    if game > 0 { 
        args = append(args, Where{Column:"game", Value:game, Compare:"="})
     }
    
    return p.Find(args)
}

func (p *GameuserManager) CountByGame(game int64, args ...interface{}) int {
    if game > 0 {
        args = append(args, Where{Column:"game", Value:game, Compare:"="})
    }
    
    return p.Count(args)
}

