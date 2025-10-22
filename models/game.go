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

type Game struct {
            
    Id                int64 `json:"id"`         
    Name                string `json:"name"`         
    Member                int `json:"member"`         
    Date                string `json:"date"` 
    
    Extra                    interface{} `json:"extra"`
}

type GameManager struct {
    Conn    *sql.DB
    Result  *sql.Result
    Index   string
}

func NewGameManager(conn *sql.DB) *GameManager {
    var item GameManager

    if conn == nil {
        item.Conn = NewConnection()
    } else {
        item.Conn = conn
    }

    item.Index = ""

    return &item
}

func (p *GameManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *GameManager) SetIndex(index string) {
    p.Index = index
}

func (p *GameManager) GetQuery() string {
    ret := ""

    str := "select g_id, g_name, g_member, g_date from game_tb "

    if p.Index == "" {
        ret = str
    } else {
        ret = str + " use index(" + p.Index + ") "
    }

    ret += "where 1=1 "
    

    return ret;
}

func (p *GameManager) GetQuerySelect() string {
    ret := ""

    str := "select count(*) from game_tb "

    if p.Index == "" {
        ret = str
    } else {
        ret = str + " use index(" + p.Index + ") "
    }

    return ret;
}

func (p *GameManager) Insert(item *Game) error {
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
        query = "insert into game_tb (g_id, g_name, g_member, g_date) values (?, ?, ?, ?)"
        res, err = p.Conn.Exec(query , item.Id, item.Name, item.Member, item.Date)
    } else {
        query = "insert into game_tb (g_name, g_member, g_date) values (?, ?, ?)"
        res, err = p.Conn.Exec(query , item.Name, item.Member, item.Date)
    }
    
    if err == nil {
        p.Result = &res
    } else {
        log.Println(err)
        p.Result = nil
    }

    return err
}
func (p *GameManager) Delete(id int64) error {
    if p.Conn == nil {
        return errors.New("Connection Error")
    }

    query := "delete from game_tb where g_id = ?"
    _, err := p.Conn.Exec(query, id)

    return err
}
func (p *GameManager) Update(item *Game) error {
    if p.Conn == nil {
        return errors.New("Connection Error")
    }

	query := "update game_tb set g_name = ?, g_member = ?, g_date = ? where g_id = ?"
	_, err := p.Conn.Exec(query , item.Name, item.Member, item.Date, item.Id)

    return err
}

func (p *GameManager) GetIdentity() int64 {
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

func (p *GameManager) ReadRow(rows *sql.Rows) *Game {
    var item Game
    var err error

    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Name, &item.Member, &item.Date)
    } else {
        return nil
    }

    if err != nil {
        return nil
    } else {

        item.Extra = map[string]interface{}{


        }
        return &item
    }
}

func (p *GameManager) ReadRows(rows *sql.Rows) *[]Game {
    var items []Game

    for rows.Next() {
        var item Game

    
        err := rows.Scan(&item.Id, &item.Name, &item.Member, &item.Date)

        if err != nil {
           log.Printf("ReadRows error : %v\n", err)
           break
        }

        item.Extra = map[string]interface{}{


        }
        items = append(items, item)
    }


     return &items
}

func (p *GameManager) Get(id int64) *Game {
    if p.Conn == nil {
        return nil
    }

    query := p.GetQuery() + " and g_id = ?"

    
    
    rows, err := p.Conn.Query(query, id)

    if err != nil {
        log.Printf("query error : %v, %v\n", err, query)
        return nil
    }

    defer rows.Close()

    return p.ReadRow(rows)
}

func (p *GameManager) Count(args []interface{}) int {
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
                query += " and g_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
            } else if item.Compare == "between" {
                query += " and g_" + item.Column + " between ? and ?"

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                query += " and g_" + item.Column + " " + item.Compare + " ?"
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

func (p *GameManager) Find(args []interface{}) *[]Game {
    if p.Conn == nil {
        var items []Game
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
                query += " and g_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
            } else if item.Compare == "between" {
                query += " and g_" + item.Column + " between ? and ?"

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                query += " and g_" + item.Column + " " + item.Compare + " ?"
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
            orderby = "g_id desc"
        } else {
            orderby = "g_" + orderby
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
            orderby = "g_id"
        } else {
            orderby = "g_" + orderby
        }
        query += " order by " + orderby
    }

    rows, err := QueryArray(p.Conn, query, params)

    if err != nil {
        log.Printf("query error : %v, %v\n", err, query)
        var items []Game
        return &items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}


