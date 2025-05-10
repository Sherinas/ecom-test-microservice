package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresPort     string
	JWTSecret        string
	AuthPort         string
	PostgresDBName   string
}

func LoadConfig() *Config {
	if err := godotenv.Load("../.env"); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file")
	}

	return &Config{
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresDBName:   os.Getenv("POSTGRESDB_NAME"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		AuthPort:         os.Getenv("AUTH_PORT"),
	}
}
