package app

import (
	"github.com/joho/godotenv"
	"log"
)

// Иннициализация переменных окружения
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func RunServiceInstance() {}
