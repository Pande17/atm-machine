package config

import (
	"log"

	"github.com/joho/godotenv"
)

func InitEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Failed loading .env file, using system env")
	}
}
