package main

import (
	"net"

	"github.com/Sherinas/ecommerce-microservices/adminservice/config"
	"github.com/Sherinas/ecommerce-microservices/adminservice/internal/client"
	"github.com/Sherinas/ecommerce-microservices/adminservice/internal/handler"
	"github.com/Sherinas/ecommerce-microservices/adminservice/logger"
	"github.com/Sherinas/ecommerce-microservices/adminservice/pb/admin"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger()

	clients, err := client.NewClients(cfg.ProductServiceAddr, cfg.UserServiceAddr, log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize clients")
	}

	lis, err := net.Listen("tcp", ":"+cfg.AdminPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	server := handler.NewAdminServer(clients, log, cfg.JWTSecret)
	grpcServer := grpc.NewServer()
	admin.RegisterAdminServiceServer(grpcServer, server)

	log.Info().Str("port", cfg.AdminPort).Msg("Starting Admin Service")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}
