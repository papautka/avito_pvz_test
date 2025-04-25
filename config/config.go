package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	DsnDb string
}

type AuthConfig struct {
	AuthTokenModerator string
	AuthTokenClient    string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return nil
	}
	return &Config{
		Db: DbConfig{
			DsnDb: os.Getenv("DSN_DB"),
		},
		Auth: AuthConfig{
			AuthTokenModerator: os.Getenv("TOKEN_MODERATOR"),
			AuthTokenClient:    os.Getenv("TOKEN_CLIENT"),
		},
	}
}
