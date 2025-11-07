package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load("configs/config.env")
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}
}
