package handler

import (
	"context"

	"github.com/Sherinas/ecommerce-microservices/Product-service/internal/repository"
	"github.com/Sherinas/ecommerce-microservices/Product-service/pb/product"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ProductServer struct {
	product.UnimplementedProductServiceServer
	repo   *repository.ProductRepository
	logger *zerolog.Logger
}

func NewProductServer(repo *repository.ProductRepository, logger *zerolog.Logger) *ProductServer {
	return &ProductServer{repo: repo, logger: logger}
}

func (s *ProductServer) AddProduct(ctx context.Context, req *product.AddProductRequest) (*product.AddProductResponse, error) {
	s.logger.Info().Str("name", req.Name).Msg("Processing AddProduct request")

	id, err := s.repo.CreateProduct(req.Name, req.Price, req.Quantity)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create product")
	}

	return &product.AddProductResponse{Id: id}, nil
}

func (s *ProductServer) UpdateProduct(ctx context.Context, req *product.UpdateProductRequest) (*product.UpdateProductResponse, error) {
	s.logger.Info().Int64("id", req.Id).Msg("Processing UpdateProduct request")

	err := s.repo.UpdateProduct(req.Id, req.Name, req.Price, req.Quantity)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to update product")
	}

	return &product.UpdateProductResponse{Success: true}, nil
}

func (s *ProductServer) DeleteProduct(ctx context.Context, req *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	s.logger.Info().Int64("id", req.Id).Msg("Processing DeleteProduct request")

	err := s.repo.DeleteProduct(req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete product")
	}

	return &product.DeleteProductResponse{Success: true}, nil
}

func (s *ProductServer) ListAllProducts(ctx context.Context, req *product.ListAllProductsRequest) (*product.ListAllProductsResponse, error) {
	s.logger.Info().Msg("Processing ListAllProducts request")

	products, err := s.repo.ListAllProducts()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list products")
	}

	protoProducts := make([]*product.Product, len(products))
	for i, p := range products {
		protoProducts[i] = &product.Product{
			Id:       int64(p.ID),
			Name:     p.Name,
			Price:    p.Price,
			Quantity: p.Quantity,
		}
	}

	return &product.ListAllProductsResponse{Products: protoProducts}, nil
}

func (s *ProductServer) GetProductById(ctx context.Context, req *product.GetProductByIdRequest) (*product.GetProductByIdResponse, error) {
	s.logger.Info().Int64("id", req.Id).Msg("Processing GetProductById request")

	p, err := s.repo.GetProductById(req.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "failed to get product")
	}

	return &product.GetProductByIdResponse{
		Id:       int64(p.ID),
		Name:     p.Name,
		Price:    p.Price,
		Quantity: p.Quantity,
	}, nil
}
