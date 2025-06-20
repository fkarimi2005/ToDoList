package db

import (
	"ToDoList/internal/configs"
	"ToDoList/logger"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
)

var db *sqlx.DB

func ConnectDB() error {
	cfg := configs.AppSettings.Postgres
	dsn := fmt.Sprintf(`host=%s
							port=%s 
							user=%s 
							password=%s 
							dbname=%s 
							sslmode=disable`, cfg.Host, cfg.Port, cfg.User, os.Getenv("DB_PASSWORD"), cfg.Database)
	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Error.Printf("[db] ConnectDB(): error during connect to postgres: %s", err.Error())
		return err
	}

	return nil
}
func CloseDB() error {
	return db.Close()
}
func GetDBConn() *sqlx.DB {
	return db
}
