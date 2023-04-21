package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // import MySQL driver
)

func ConnectDB(dsn string) (*sql.DB, error) {
	// Open a database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
