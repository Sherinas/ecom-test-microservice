package main

import (
	"net"

	"github.com/Sherinas/ecommerce-microservices/Product-service/config"
	"github.com/Sherinas/ecommerce-microservices/Product-service/internal/db"
	"github.com/Sherinas/ecommerce-microservices/Product-service/internal/handler"
	"github.com/Sherinas/ecommerce-microservices/Product-service/internal/repository"
	"github.com/Sherinas/ecommerce-microservices/Product-service/logger"
	"github.com/Sherinas/ecommerce-microservices/Product-service/pb/product"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	log := logger.NewLogger()

	db, err := db.NewDB(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	repo := repository.NewProductRepository(db, log)
	server := handler.NewProductServer(repo, log)

	lis, err := net.Listen("tcp", ":"+cfg.ProductPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	grpcServer := grpc.NewServer()
	product.RegisterProductServiceServer(grpcServer, server)

	log.Info().Str("port", cfg.ProductPort).Msg("Starting Product Service")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("Failed to serve")
	}
}
