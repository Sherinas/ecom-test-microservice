package repository

import (
	"github.com/Sherinas/ecommerce-microservices/product-service/internal/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewProductRepository(db *gorm.DB, logger *zerolog.Logger) *ProductRepository {
	return &ProductRepository{db: db, logger: logger}
}

func (r *ProductRepository) CreateProduct(name string, price float32, quantity int32) (int64, error) {
	product := &models.Product{
		Name:     name,
		Price:    price,
		Quantity: quantity,
	}

	if err := r.db.Create(product).Error; err != nil {
		r.logger.Error().Err(err).Str("name", name).Msg("Failed to create product")
		return 0, err
	}

	return int64(product.ID), nil
}

func (r *ProductRepository) UpdateProduct(id int64, name string, price float32, quantity int32) error {
	result := r.db.Model(&models.Product{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":     name,
		"price":    price,
		"quantity": quantity,
	})
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Int64("id", id).Msg("Failed to update product")
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.logger.Warn().Int64("id", id).Msg("Product not found")
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(id int64) error {
	result := r.db.Where("id = ?", id).Delete(&models.Product{})
	if result.Error != nil {
		r.logger.Error().Err(result.Error).Int64("id", id).Msg("Failed to delete product")
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.logger.Warn().Int64("id", id).Msg("Product not found")
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ProductRepository) ListAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Find(&products).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to list products")
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetProductById(id int64) (*models.Product, error) {
	var product models.Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		r.logger.Error().Err(err).Int64("id", id).Msg("Failed to get product")
		return nil, err
	}
	return &product, nil
}
