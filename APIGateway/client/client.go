package client

import (
	"github.com/Sherinas/ecommerce-microservices/APIGateway/pb/admin"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/pb/auth"
	"github.com/Sherinas/ecommerce-microservices/APIGateway/pb/product"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Clients struct {
	AuthClient    auth.AuthServiceClient
	AdminClient   admin.AdminServiceClient
	ProductClient product.ProductServiceClient
}

func NewClients(authAddr, adminAddr, productAddr string, logger *zerolog.Logger) (*Clients, error) {
	authConn, err := grpc.Dial(authAddr, grpc.WithInsecure())
	if err != nil {
		logger.Error().Err(err).Str("addr", authAddr).Msg("Failed to connect to Auth Service")
		return nil, err
	}

	adminConn, err := grpc.Dial(adminAddr, grpc.WithInsecure())
	if err != nil {
		logger.Error().Err(err).Str("addr", adminAddr).Msg("Failed to connect to Admin Service")
		return nil, err
	}

	productConn, err := grpc.Dial(productAddr, grpc.WithInsecure())
	if err != nil {
		logger.Error().Err(err).Str("addr", productAddr).Msg("Failed to connect to Product Service")
		return nil, err
	}

	return &Clients{
		AuthClient:    auth.NewAuthServiceClient(authConn),
		AdminClient:   admin.NewAdminServiceClient(adminConn),
		ProductClient: product.NewProductServiceClient(productConn),
	}, nil
}
