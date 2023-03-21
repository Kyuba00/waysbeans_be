package repositories

import (
	"waysbeans_be/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindProducts() ([]models.Product, error)
	GetProduct(ID int) (models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
	UpdateProduct(product models.Product) (models.Product, error)
	DeleteProduct(product models.Product) (models.Product, error)
}

// PRODUCT REPOSITORY FUNCTION
func RepositoryProduct(db *gorm.DB) *repository {
	return &repository{db}
}

// FIND PRODUCTS FROM DATABASE
func (r *repository) FindProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error

	return products, err
}

// FIND USER FROM DATABASE
func (r *repository) GetProduct(ID int) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, ID).Error
	return product, err
}

// CREATE PRODUCT TO DATABASE
func (r *repository) CreateProduct(product models.Product) (models.Product, error) {
	err := r.db.Create(&product).Error
	return product, err
}

// UPDATE PRODUCT TO DATABASE
func (r *repository) UpdateProduct(product models.Product) (models.Product, error) {
	err := r.db.Save(&product).Error
	return product, err
}

// DELETE PRODUCT FROM DATABASE
func (r *repository) DeleteProduct(product models.Product) (models.Product, error) {
	err := r.db.Delete(&product).Error
	return product, err
}
