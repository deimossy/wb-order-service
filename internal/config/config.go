package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	AppPort     string
	AppHost     string
	once        sync.Once
)

func LoadConfig() {
	once.Do(func() {
		err := godotenv.Load(".env")
		if err != nil {
			log.Println("Error loading .env file")
		}
		AppPort = os.Getenv("APP_PORT")
		AppHost = os.Getenv("APP_HOST")
	})
}