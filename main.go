package main

import (
	_ "ToDoList/docs"
	"ToDoList/internal/configs"
	"ToDoList/internal/controller"
	"ToDoList/internal/db"
	"ToDoList/logger"
	"fmt"
	"log"
	"os"
)

// @title           ToDoList API
// @version         1.0
// @description     API для управления задачами с авторизацией
// @host            localhost:8089
// @BasePath        /
// @schemes         http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	fmt.Println("DB password from env:", os.Getenv("DB_PASSWORD"))

	if err := configs.ReadSettings(); err != nil {
		log.Fatal("Error during reading settings: ", err)
	}

	if err := logger.Init(); err != nil {
		log.Println("Error during init logger: ", err)
	}

	if err := db.ConnectDB(); err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	if err := db.InitMigrations(); err != nil {
		log.Fatal("Error initializing database: ", err)
	}

	if err := controller.RunServer(); err != nil {
		log.Println("Error during HTTP server: ", err)
	}
}
