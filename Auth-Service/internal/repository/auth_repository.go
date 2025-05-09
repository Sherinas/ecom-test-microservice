package repository

import (
	"github.com/Sherinas/ecommerce-microservices/Auth-Service/internal/models"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewAuthRepository(db *gorm.DB, logger *zerolog.Logger) *AuthRepository {
	return &AuthRepository{db: db, logger: logger}
}

func (r *AuthRepository) CreateUser(email, password, role string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to hash password")
		return nil, err
	}

	user := &models.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := r.db.Create(user).Error; err != nil {
		r.logger.Error().Err(err).Str("email", email).Msg("Failed to create user")
		return nil, err
	}

	return user, nil
}

func (r *AuthRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		r.logger.Error().Err(err).Str("email", email).Msg("Failed to find user")
		return nil, err
	}
	return &user, nil
}
