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

type User struct {
            
    Id                int64 `json:"id"`         
    Loginid                string `json:"loginid"`         
    Passwd                string `json:"passwd"`         
    Name                string `json:"name"`         
    Date                string `json:"date"` 
    
    Extra                    interface{} `json:"extra"`
}

type UserManager struct {
    Conn    *sql.DB
    Result  *sql.Result
    Index   string
}

func NewUserManager(conn *sql.DB) *UserManager {
    var item UserManager

    if conn == nil {
        item.Conn = NewConnection()
    } else {
        item.Conn = conn
    }

    item.Index = ""

    return &item
}

func (p *UserManager) Close() {
    if p.Conn != nil {
        p.Conn.Close()
    }
}

func (p *UserManager) SetIndex(index string) {
    p.Index = index
}

func (p *UserManager) GetQuery() string {
    ret := ""

    str := "select u_id, u_loginid, u_passwd, u_name, u_date from user_tb "

    if p.Index == "" {
        ret = str
    } else {
        ret = str + " use index(" + p.Index + ") "
    }

    ret += "where 1=1 "
    

    return ret;
}

func (p *UserManager) GetQuerySelect() string {
    ret := ""

    str := "select count(*) from user_tb "

    if p.Index == "" {
        ret = str
    } else {
        ret = str + " use index(" + p.Index + ") "
    }

    return ret;
}

func (p *UserManager) Insert(item *User) error {
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
        query = "insert into user_tb (u_id, u_loginid, u_passwd, u_name, u_date) values (?, ?, ?, ?, ?)"
        res, err = p.Conn.Exec(query , item.Id, item.Loginid, item.Passwd, item.Name, item.Date)
    } else {
        query = "insert into user_tb (u_loginid, u_passwd, u_name, u_date) values (?, ?, ?, ?)"
        res, err = p.Conn.Exec(query , item.Loginid, item.Passwd, item.Name, item.Date)
    }
    
    if err == nil {
        p.Result = &res
    } else {
        log.Println(err)
        p.Result = nil
    }

    return err
}
func (p *UserManager) Delete(id int64) error {
    if p.Conn == nil {
        return errors.New("Connection Error")
    }

    query := "delete from user_tb where u_id = ?"
    _, err := p.Conn.Exec(query, id)

    return err
}
func (p *UserManager) Update(item *User) error {
    if p.Conn == nil {
        return errors.New("Connection Error")
    }

	query := "update user_tb set u_loginid = ?, u_passwd = ?, u_name = ?, u_date = ? where u_id = ?"
	_, err := p.Conn.Exec(query , item.Loginid, item.Passwd, item.Name, item.Date, item.Id)

    return err
}

func (p *UserManager) GetIdentity() int64 {
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

func (p *UserManager) ReadRow(rows *sql.Rows) *User {
    var item User
    var err error

    

    if rows.Next() {
        err = rows.Scan(&item.Id, &item.Loginid, &item.Passwd, &item.Name, &item.Date)
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

func (p *UserManager) ReadRows(rows *sql.Rows) *[]User {
    var items []User

    for rows.Next() {
        var item User

    
        err := rows.Scan(&item.Id, &item.Loginid, &item.Passwd, &item.Name, &item.Date)

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

func (p *UserManager) Get(id int64) *User {
    if p.Conn == nil {
        return nil
    }

    query := p.GetQuery() + " and u_id = ?"

    
    
    rows, err := p.Conn.Query(query, id)

    if err != nil {
        log.Printf("query error : %v, %v\n", err, query)
        return nil
    }

    defer rows.Close()

    return p.ReadRow(rows)
}

func (p *UserManager) Count(args []interface{}) int {
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
                query += " and u_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
            } else if item.Compare == "between" {
                query += " and u_" + item.Column + " between ? and ?"

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                query += " and u_" + item.Column + " " + item.Compare + " ?"
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

func (p *UserManager) Find(args []interface{}) *[]User {
    if p.Conn == nil {
        var items []User
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
                query += " and u_id in (" + strings.Trim(strings.Replace(fmt.Sprint(item.Value), " ", ", ", -1), "[]") + ")"
            } else if item.Compare == "between" {
                query += " and u_" + item.Column + " between ? and ?"

                s := item.Value.([2]string)
                params = append(params, s[0])
                params = append(params, s[1])
            } else {
                query += " and u_" + item.Column + " " + item.Compare + " ?"
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
            orderby = "u_id desc"
        } else {
            orderby = "u_" + orderby
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
            orderby = "u_id"
        } else {
            orderby = "u_" + orderby
        }
        query += " order by " + orderby
    }

    rows, err := QueryArray(p.Conn, query, params)

    if err != nil {
        log.Printf("query error : %v, %v\n", err, query)
        var items []User
        return &items
    }

    defer rows.Close()

    return p.ReadRows(rows)
}


