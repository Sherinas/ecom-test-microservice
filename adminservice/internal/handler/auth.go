package handler

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (s *AdminServer) validateAdmin(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		s.logger.Error().Msg("No metadata provided")
		return status.Error(codes.Unauthenticated, "no metadata provided")
	}

	tokenStrings, ok := md["authorization"]
	if !ok || len(tokenStrings) == 0 {
		s.logger.Error().Msg("No authorization token provided")
		return status.Error(codes.Unauthenticated, "no authorization token")
	}

	tokenString := tokenStrings[0]
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		s.logger.Error().Err(err).Msg("Invalid JWT token")
		return status.Error(codes.Unauthenticated, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error().Msg("Invalid token claims")
		return status.Error(codes.Unauthenticated, "invalid token claims")
	}

	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		s.logger.Error().Str("role", role).Msg("User is not an admin")
		return status.Error(codes.PermissionDenied, "admin access required")
	}

	return nil
}
