package models

import (
	"mighty/config"

	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type PagingType struct {
	Page     int
	Pagesize int
}

type OrderingType struct {
	Order string
}

type LimitType struct {
	Limit int
}

type Option struct {
	Page     int
	Pagesize int
	Order    string
	Limit    int
}

type Where struct {
	Column  string
	Value   interface{}
	Compare string
}

func Paging(page int, pagesize int) PagingType {
	return PagingType{Page: page, Pagesize: pagesize}
}

func Ordering(order string) OrderingType {
	return OrderingType{Order: order}
}

func Limit(limit int) LimitType {
	return LimitType{Limit: limit}
}

func GetConnection() *sql.DB {
	r1, err := sql.Open(config.Database, config.ConnectionString)
	if err != nil {
		log.Println("Database Connect Error")
		return nil
	}

	r1.SetMaxOpenConns(100)
	r1.SetMaxIdleConns(10)
	r1.SetConnMaxLifetime(5 * time.Minute)

	return r1
}

func NewConnection() *sql.DB {
	db := GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(100 * time.Millisecond)

	db = GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(500 * time.Millisecond)

	db = GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(1 * time.Second)

	db = GetConnection()

	if db != nil {
		return db
	}

	time.Sleep(2 * time.Second)

	db = GetConnection()

	return db
}

func QueryArray(db *sql.DB, query string, items []interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error

	size := len(items)

	if size == 0 {
		rows, err = db.Query(query)
	} else if size == 1 {
		rows, err = db.Query(query, items[0])
	} else if size == 2 {
		rows, err = db.Query(query, items[0], items[1])
	} else if size == 3 {
		rows, err = db.Query(query, items[0], items[1], items[2])
	} else if size == 4 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3])
	} else if size == 5 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4])
	} else if size == 6 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4], items[5])
	} else if size == 7 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6])
	} else if size == 8 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7])
	} else if size == 9 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8])
	} else if size == 10 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9])
	} else if size == 11 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10])
	} else if size == 12 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10], items[11])
	} else if size == 13 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10], items[11], items[12])
	} else if size == 14 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10], items[11], items[12], items[13])
	} else if size == 15 {
		rows, err = db.Query(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10], items[11], items[12], items[13], items[14])
	}

	return rows, err
}

func ExecArray(db *sql.DB, query string, items []interface{}) error {
	var err error

	size := len(items)

	if size == 0 {
		_, err = db.Exec(query)
	} else if size == 1 {
		_, err = db.Exec(query, items[0])
	} else if size == 2 {
		_, err = db.Exec(query, items[0], items[1])
	} else if size == 3 {
		_, err = db.Exec(query, items[0], items[1], items[2])
	} else if size == 4 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3])
	} else if size == 5 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4])
	} else if size == 6 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4], items[5])
	} else if size == 7 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6])
	} else if size == 8 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7])
	} else if size == 9 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8])
	} else if size == 10 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9])
	} else if size == 11 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10])
	} else if size == 12 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10], items[11])
	} else if size == 13 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10], items[11], items[12])
	} else if size == 14 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10], items[11], items[12], items[13])
	} else if size == 15 {
		_, err = db.Exec(query, items[0], items[1], items[2], items[3], items[4], items[5], items[6], items[7], items[8], items[9], items[10], items[11], items[12], items[13], items[14])
	}

	return err
}

func InitDate() string {
	return "1000-01-01 00:00:00"
}
