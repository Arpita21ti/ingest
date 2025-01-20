package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	// Try to load the .env file but don’t fail if it doesn’t exist
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, loading environment variables from the system")
	}
	return nil
}