package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	AuthServiceAddr    string
	AdminServiceAddr   string
	ProductServiceAddr string
	JWTSecret          string
	AdminSecret        string
	APIGatewayPort     string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Warn().Err(err).Msg("Error loading .env file")
	}

	return &Config{
		AuthServiceAddr:    os.Getenv("AUTH_SERVICE_ADDR"),
		AdminServiceAddr:   os.Getenv("ADMIN_SERVICE_ADDR"),
		ProductServiceAddr: os.Getenv("PRODUCT_SERVICE_ADDR"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		AdminSecret:        os.Getenv("ADMIN_SECRET"),
		APIGatewayPort:     os.Getenv("API_GATEWAY_PORT"),
	}
}
