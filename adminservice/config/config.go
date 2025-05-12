package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ProductServiceAddr string
	UserServiceAddr    string
	JWTSecret          string
	AdminPort          string
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file")
	}

	return &Config{
		ProductServiceAddr: os.Getenv("PRODUCT_SERVICE_ADDR"),
		UserServiceAddr:    os.Getenv("USER_SERVICE_ADDR"),
		AdminPort:          os.Getenv("ADMIN_PORT"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}
}
