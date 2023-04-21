package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func GetDB() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/cobra")
		if err != nil {
			panic(err)
		}
	}

	return db
}
