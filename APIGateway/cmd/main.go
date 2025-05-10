package main

import (
	"github.com/Sherinas/ecommerce-microservices/APIGateway/client"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/config"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/internal/handler"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/internal/middleware"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/internal/util"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger()

	clients, err := client.NewClients(cfg.AuthServiceAddr, cfg.AdminServiceAddr, cfg.ProductServiceAddr, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize gRPC clients")
	}

	validator := util.NewAuthValidator(log, cfg.JWTSecret)
	handler := handler.NewGatewayHandler(clients, log, validator, cfg.AdminSecret)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggingMiddleware(log))
	r.Use(middleware.ErrorHandlingMiddleware())

	handler.RegisterRoutes(r)

	log.Info().Str("port", cfg.APIGatewayPort).Msg("Starting API Gateway")
	if err := r.Run(":" + cfg.APIGatewayPort); err != nil {
		log.Fatal().Err(err).Msg("Failed to start API Gateway")
	}
}
