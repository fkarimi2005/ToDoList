package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
)

func ReadSettings() {
	fmt.Println("Loading file .env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Reading confi)

}
