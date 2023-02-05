package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func NewEnvConfig(keys ...string) []string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("error loading .env file %v\n", err)
	}

	keyEnv := []string{}

	for _, key := range keys {
		env := os.Getenv(key)
		keyEnv = append(keyEnv, env)
	}

	return keyEnv
}
