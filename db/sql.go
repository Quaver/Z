package db

import (
	"example.com/Quaver/Z/config"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

type RowScanner interface {
	Scan(dest ...interface{}) error
}

var SQL *sqlx.DB

// InitializeSQL Initializes the SQL database connection
func InitializeSQL() {
	if SQL != nil {
		return
	}

	credentials := config.Instance.SQL
	connStr := fmt.Sprintf("%v:%v@tcp(%v)/%v", credentials.Username, credentials.Password, credentials.Host, credentials.Database)

	db, err := sqlx.Connect("mysql", connStr)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatalln(err)
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
