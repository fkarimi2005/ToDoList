package configs

import (
	"ToDoList/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var AppSettings model.Configs

func ReadSettings() error {
	fmt.Println("Loading file .env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Reading Settings file: configs/configs.json ")
	configFile, err := os.Open("internal/configs/configs.json")
	if err != nil {
		fmt.Println("Error during opening configs.json file")
	}
	defer func(configFile *os.File) {
		if err := configFile.Close(); err != nil {
			log.Fatal("Error closing configs.json file", err.Error())
		}
	}(configFile)
	if err := json.NewDecoder(configFile).Decode(&AppSettings); err != nil {
		return errors.New(fmt.Sprintf("Error decoding configs.json file %s", err.Error()))
	}
	return nil
}
