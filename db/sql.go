package db

import (
	"database/sql"
	"example.com/Quaver/Z/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type RowScanner interface {
	Scan(dest ...interface{}) error
}

var SQL *sql.DB

// InitializeSQL Initializes the SQL database connection
func InitializeSQL() {
	if SQL != nil {
		return
	}

	credentials := config.Instance.SQL
	connStr := fmt.Sprintf("%v:%v@tcp(%v)/%v", credentials.Username, credentials.Password, credentials.Host, credentials.Database)

	db, err := sql.Open("mysql", connStr)

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	SQL = db
	log.Println("Successfully connected to SQL database")
}

// CloseSQLConnection Closes the existing SQL connection
func CloseSQLConnection() {
	if SQL == nil {
		return
	}

	err := SQL.Close()

	if err != nil {
		return
	}

	log.Println("SQL database connection has been closed")
	SQL = nil
}
