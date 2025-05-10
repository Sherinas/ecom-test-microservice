package handler

import (
	"context"

	jwtauth "github.com/Sherinas/ecommerce-microservices/auth-service/internal/jwt"
	"github.com/Sherinas/ecommerce-microservices/auth-service/internal/repository"
	"github.com/Sherinas/ecommerce-microservices/auth-service/internal/utils"
	"github.com/Sherinas/ecommerce-microservices/auth-service/pb/auth"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gorm.io/gorm"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	repo   *repository.AuthRepository
	logger *zerolog.Logger
	secret string
}

func NewAuthServer(db *gorm.DB, logger *zerolog.Logger, secret string) *AuthServer {
	repo := repository.NewAuthRepository(db, logger)
	return &AuthServer{repo: repo, logger: logger, secret: secret}
}

func (s *AuthServer) SignUp(ctx context.Context, req *auth.SignUpRequest) (*auth.SignUpResponse, error) {
	s.logger.Info().Str("email", req.Email).Msg("Processing SignUp request")

	// Check if user already exists
	if _, err := s.repo.FindUserByEmail(req.Email); err == nil {
		s.logger.Warn().Str("email", req.Email).Msg("User already exists")
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	// Create user
	user, err := s.repo.CreateUser(req.Email, req.Password, req.Role)
	if err != nil {
		s.logger.Error().Err(err).Str("email", req.Email).Msg("Failed to create user")
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	// Generate JWT
	token, err := jwtauth.GenerateJWT(user.Email, s.secret, user.Role)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate JWT")
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &auth.SignUpResponse{Token: token}, nil
}

func (s *AuthServer) SignIn(ctx context.Context, req *auth.SignInRequest) (*auth.SignInResponse, error) {
	s.logger.Info().Str("email", req.Email).Msg("Processing SignIn request")

	// Find user
	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		s.logger.Error().Err(err).Str("email", req.Email).Msg("User not found")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Verify password
	if err := utils.VerifyPassword(user.Password, req.Password); err != nil {
		s.logger.Error().Err(err).Str("email", req.Email).Msg("Invalid password")
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Generate JWT
	token, err := jwtauth.GenerateJWT(user.Email, s.secret, user.Role)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate JWT")
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	return &auth.SignInResponse{Token: token}, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	s.logger.Info().Msg("Processing Logout request")
	// In production, add token to blacklist (e.g., Redis)
	return &auth.LogoutResponse{Success: true}, nil
}
