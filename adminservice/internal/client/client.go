package client

import (
	"github.com/Sherinas/ecommerce-microservices/adminservice/pb/auth"
	"github.com/Sherinas/ecommerce-microservices/adminservice/pb/product"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Clients struct {
	ProductClient product.ProductServiceClient
	AuthClient    auth.AuthServiceClient
}

func NewClients(productAddr, userAddr string, logger *zerolog.Logger) (*Clients, error) {
	productConn, err := grpc.Dial(productAddr, grpc.WithInsecure())
	if err != nil {
		logger.Error().Err(err).Str("addr", productAddr).Msg("Failed to connect to Product Service")
		return nil, err
	}

	userConn, err := grpc.Dial(userAddr, grpc.WithInsecure())
	if err != nil {
		logger.Error().Err(err).Str("addr", userAddr).Msg("Failed to connect to User Service")
		return nil, err
	}

	return &Clients{
		ProductClient: product.NewProductServiceClient(productConn),
		AuthClient:    auth.NewAuthServiceClient(userConn),
	}, nil
}
