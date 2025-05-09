package main

import (
	"net"

	"github.com/Sherinas/ecommerce-microservices/Auth-Service/config"
	"github.com/Sherinas/ecommerce-microservices/Auth-Service/internal/db"
	"github.com/Sherinas/ecommerce-microservices/Auth-Service/internal/handler"
	"github.com/Sherinas/ecommerce-microservices/Auth-Service/logger"
	"github.com/Sherinas/ecommerce-microservices/Auth-Service/pb/auth"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger()

	db, err := db.NewDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	lis, err := net.Listen("tcp", ":"+cfg.AuthPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	server := handler.NewAuthServer(db, log, cfg.JWTSecret)
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, server)

	log.Info().Str("port", cfg.AuthPort).Msg("Starting Auth Service")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}
