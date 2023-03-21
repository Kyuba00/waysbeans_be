package repositories

import (
	"waysbeans_be/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	// CART
	FindCartByTransactionID(TransactionID int) ([]models.Cart, error)
	GetCart(CartID int) (models.Cart, error)
	CreateCart(cart models.Cart) (models.Cart, error)
	DeleteCart(cart models.Cart) (models.Cart, error)

	// TRANSACTION
	GetCartTransactionByUserID(UserID int) (models.Transaction, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)

	// GET PRODUCT ID
	GetProductCartByID(ProductID int) (models.Product, error)
}

// PRODUCT REPOSITORY FUNCTION
func RepositoryCart(db *gorm.DB) *repository {
	return &repository{db}
}

// FIND CART BY TRANSACTION ID FROM DATABASE
func (r *repository) FindCartByTransactionID(TransactionID int) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Preload("Product").Find(&carts, "transaction_id = ?", TransactionID).Error
	return carts, err
}

// GET CART BY ID FROM DATABASE
func (r *repository) GetCart(CartID int) (models.Cart, error) {
	var cart models.Cart
	err := r.db.Preload("Product").First(&cart, CartID).Error
	return cart, err
}

// CREATE TRANSACTION TO DATABASE
func (r *repository) CreateCart(cart models.Cart) (models.Cart, error) {
	err := r.db.Create(&cart).Error
	return cart, err
}

// DELETE TRANSACTION FROM DATABASE
func (r *repository) DeleteCart(cart models.Cart) (models.Cart, error) {
	err := r.db.Delete(&cart).Error
	return cart, err
}

// GET TRANSACTION BY USER ID FOR CART FROM DATABASE
func (r *repository) GetCartTransactionByUserID(UserID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Product").Where("user_id = ?", UserID).Where("status = ?", "waiting").Find(&transaction).Error
	return transaction, err
}

// CREATE TRANSACTION TO DATABASE
func (r *repository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Product").Create(&transaction).Error
	return transaction, err
}

// GET PRODUCT ORDER BY ID FROM DATABASE
func (r *repository) GetProductCartByID(ProductID int) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, ProductID).Error
	return product, err
}
