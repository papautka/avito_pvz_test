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
	}

	dsn := os.Getenv("DSN_DB")
	if dsn == "" {
		log.Fatal("DSN_DB is not set in environment")
	}
	return &Config{
		Db: DbConfig{
			DsnDb: dsn,
		},
		Auth: AuthConfig{
			AuthTokenModerator: os.Getenv("TOKEN_MODERATOR"),
			AuthTokenClient:    os.Getenv("TOKEN_CLIENT"),
		},
	}
}
