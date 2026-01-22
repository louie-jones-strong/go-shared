package env

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Printf(".env file not found: %v", err)
	}
	return nil
}

func GetKey(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("%v environment variable not set. Please add it to your .env file", key)
	}
	return val, nil
}
