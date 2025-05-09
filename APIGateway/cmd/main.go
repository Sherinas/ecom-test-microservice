package main

import (
	"github.com/Sherinas/ecommerce-microservices/APIGateway/internal/config"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/internal/gateway"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger()

	clients, err := gateway.NewClients(cfg.AuthServiceAddr, cfg.AdminServiceAddr, cfg.ProductServiceAddr, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize gRPC clients")
	}

	validator := gateway.NewAuthValidator(log, cfg.JWTSecret)
	handler := gateway.NewGatewayHandler(clients, log, validator, cfg.AdminSecret)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gateway.LoggingMiddleware(log))
	r.Use(gateway.ErrorHandlingMiddleware())

	handler.RegisterRoutes(r)

	log.Info().Str("port", cfg.APIGatewayPort).Msg("Starting API Gateway")
	if err := r.Run(":" + cfg.APIGatewayPort); err != nil {
		log.Fatal().Err(err).Msg("Failed to start API Gateway")
	}
}
