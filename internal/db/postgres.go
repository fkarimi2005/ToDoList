package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func ConnectDB() error {
	dsn := "host=127.0.0.1 port=5432 user=postgres password=Firuzshoh2005 dbname=todolist_db sslmode=disable"
	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}

	fmt.Println("Connected to PostgreSQL!")
	return nil
}
func CloseDB() error {
	return db.Close()
}
func GetDBConn() *sqlx.DB {
	return db
}
