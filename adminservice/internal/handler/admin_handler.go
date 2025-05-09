package handler

import (
	"context"

	"github.com/Sherinas/ecommerce-microservices/adminservice/internal/client"
	"github.com/Sherinas/ecommerce-microservices/adminservice/pb/admin"
	product "github.com/Sherinas/ecommerce-microservices/adminservice/pb/product"
	user "github.com/Sherinas/ecommerce-microservices/adminservice/pb/user"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminServer struct {
	admin.UnimplementedAdminServiceServer
	clients *client.Clients
	logger  *zerolog.Logger
	secret  string
}

func NewAdminServer(clients *client.Clients, logger *zerolog.Logger, secret string) *AdminServer {
	return &AdminServer{clients: clients, logger: logger, secret: secret}
}

func (s *AdminServer) AddProduct(ctx context.Context, req *admin.AddProductRequest) (*admin.AddProductResponse, error) {
	if err := s.validateAdmin(ctx); err != nil {
		return nil, err
	}

	s.logger.Info().Str("name", req.Name).Msg("Processing AddProduct request")

	resp, err := s.clients.ProductClient.AddProduct(ctx, &product.AddProductRequest{
		Name:     req.Name,
		Price:    req.Price,
		Quantity: req.Quantity,
	})
	if err != nil {
		s.logger.Error().Err(err).Str("name", req.Name).Msg("Failed to add product")
		return nil, status.Error(codes.Internal, "failed to add product")
	}

	return &admin.AddProductResponse{Id: resp.Id}, nil
}

func (s *AdminServer) UpdateProduct(ctx context.Context, req *admin.UpdateProductRequest) (*admin.UpdateProductResponse, error) {
	if err := s.validateAdmin(ctx); err != nil {
		return nil, err
	}

	s.logger.Info().Int64("id", req.Id).Msg("Processing UpdateProduct request")

	_, err := s.clients.ProductClient.UpdateProduct(ctx, &product.UpdateProductRequest{
		Id:       req.Id,
		Name:     req.Name,
		Price:    req.Price,
		Quantity: req.Quantity,
	})
	if err != nil {
		s.logger.Error().Err(err).Int64("id", req.Id).Msg("Failed to update product")
		return nil, status.Error(codes.Internal, "failed to update product")
	}

	return &admin.UpdateProductResponse{Success: true}, nil
}

func (s *AdminServer) DeleteProduct(ctx context.Context, req *admin.DeleteProductRequest) (*admin.DeleteProductResponse, error) {
	if err := s.validateAdmin(ctx); err != nil {
		return nil, err
	}

	s.logger.Info().Int64("id", req.Id).Msg("Processing DeleteProduct request")

	_, err := s.clients.ProductClient.DeleteProduct(ctx, &product.DeleteProductRequest{Id: req.Id})
	if err != nil {
		s.logger.Error().Err(err).Int64("id", req.Id).Msg("Failed to delete product")
		return nil, status.Error(codes.Internal, "failed to delete product")
	}

	return &admin.DeleteProductResponse{Success: true}, nil
}

func (s *AdminServer) ListAllProducts(ctx context.Context, req *admin.ListAllProductsRequest) (*admin.ListAllProductsResponse, error) {
	if err := s.validateAdmin(ctx); err != nil {
		return nil, err
	}

	s.logger.Info().Msg("Processing ListAllProducts request")

	resp, err := s.clients.ProductClient.ListAllProducts(ctx, &product.ListAllProductsRequest{})
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list products")
		return nil, status.Error(codes.Internal, "failed to list products")
	}

	products := make([]*admin.Product, len(resp.Products))
	for i, p := range resp.Products {
		products[i] = &admin.Product{
			Id:       p.Id,
			Name:     p.Name,
			Price:    p.Price,
			Quantity: p.Quantity,
		}
	}

	return &admin.ListAllProductsResponse{Products: products}, nil
}

func (s *AdminServer) GetProductById(ctx context.Context, req *admin.GetProductByIdRequest) (*admin.GetProductByIdResponse, error) {
	if err := s.validateAdmin(ctx); err != nil {
		return nil, err
	}

	s.logger.Info().Int64("id", req.Id).Msg("Processing GetProductById request")

	resp, err := s.clients.ProductClient.GetProductById(ctx, &product.GetProductByIdRequest{Id: req.Id})
	if err != nil {
		s.logger.Error().Err(err).Int64("id", req.Id).Msg("Failed to get product")
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &admin.GetProductByIdResponse{
		Id:       resp.Id,
		Name:     resp.Name,
		Price:    resp.Price,
		Quantity: resp.Quantity,
	}, nil
}

func (s *AdminServer) ListUsers(ctx context.Context, req *admin.ListUsersRequest) (*admin.ListUsersResponse, error) {
	if err := s.validateAdmin(ctx); err != nil {
		return nil, err
	}

	s.logger.Info().Msg("Processing ListUsers request")

	resp, err := s.clients.UserClient.ListUsers(ctx, &user.ListUsersRequest{})
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list users")
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	users := make([]*admin.User, len(resp.Users))
	for i, u := range resp.Users {
		users[i] = &admin.User{
			Id:        u.Id,
			Email:     u.Email,
			IsBlocked: u.IsBlocked,
		}
	}

	return &admin.ListUsersResponse{Users: users}, nil
}

func (s *AdminServer) DeleteUser(ctx context.Context, req *admin.DeleteUserRequest) (*admin.DeleteUserResponse, error) {
	if err := s.validateAdmin(ctx); err != nil {
		return nil, err
	}

	s.logger.Info().Int64("id", req.Id).Msg("Processing DeleteUser request")

	_, err := s.clients.UserClient.DeleteUser(ctx, &user.DeleteUserRequest{Id: req.Id})
	if err != nil {
		s.logger.Error().Err(err).Int64("id", req.Id).Msg("Failed to delete user")
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	return &admin.DeleteUserResponse{Success: true}, nil
}

func (s *AdminServer) BlockUser(ctx context.Context, req *admin.BlockUserRequest) (*admin.BlockUserResponse, error) {
	if err := s.validateAdmin(ctx); err != nil {
		return nil, err
	}

	s.logger.Info().Int64("id", req.Id).Msg("Processing BlockUser request")

	_, err := s.clients.UserClient.BlockUser(ctx, &user.BlockUserRequest{Id: req.Id})
	if err != nil {
		s.logger.Error().Err(err).Int64("id", req.Id).Msg("Failed to block user")
		return nil, status.Error(codes.Internal, "failed to block user")
	}

	return &admin.BlockUserResponse{Success: true}, nil
}
