package gateway

import (
	"context"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog"
)

type AuthValidator struct {
	logger *zerolog.Logger
	secret string
}

func NewAuthValidator(logger *zerolog.Logger, secret string) *AuthValidator {
	return &AuthValidator{logger: logger, secret: secret}
}

func (v *AuthValidator) ValidateAdminJWT(tokenString string) error {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(v.secret), nil
	})
	if err != nil || !token.Valid {
		v.logger.Error().Err(err).Msg("Invalid JWT token")
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		v.logger.Error().Msg("Invalid token claims")
		return jwt.ErrInvalidKey
	}

	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		v.logger.Error().Str("role", role).Msg("User is not an admin")
		return jwt.ErrInvalidKey
	}

	return nil
}

// Create gRPC metadata with JWT for downstream services
func (v *AuthValidator) CreateGRPCMetadata(tokenString string) context.Context {
	md := metadata.Pairs("authorization", "Bearer "+tokenString)
	return metadata.NewOutgoingContext(context.Background(), md)
}
