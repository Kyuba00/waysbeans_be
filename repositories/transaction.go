package repositories

import (
	"waysbeans_be/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindTransactions() ([]models.Transaction, error)
	GetTransactionByUserID(UserID int) (models.Transaction, error)
	GetUserTransactionByUserID(UserID int) ([]models.Transaction, error)
	GetTransactionNotification(ID int) (models.Transaction, error)
	GetTransactionMidtrans(ID string) (models.Transaction, error)
	UpdateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateTransactionMidtrans(status string, ID int) error
}

// PRODUCT REPOSITORY FUNCTION
func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

// FIND TRANSACTIONS IS NOT WAITING FROM DATABASE
func (r *repository) FindTransactions() ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Product").Not("status = ?", "waiting").Order("id DESC").Find(&transactions).Error
	return transactions, err
}

// GET TRANSACTION BY USER ID FROM DATABASE WHERE IS WAITING
func (r *repository) GetTransactionByUserID(UserID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Product").Where("user_id = ?", UserID).Where("status = ?", "waiting").Find(&transaction).Error
	return transaction, err
}

// GET TRANSACTION BY USER ID FROM DATABASE WHERE NOT WAITING
func (r *repository) GetUserTransactionByUserID(UserID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Product").Where("user_id = ?", UserID).Not("status = ?", "waiting").Order("id DESC").Find(&transactions).Error
	return transactions, err
}

// GET TRANSACTION NOTIFICATION FROM MIDTRANS
func (r *repository) GetTransactionNotification(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.First(&transaction, ID).Error
	return transaction, err
}

// GET TRANSACTION FROM MIDTRANS
func (r *repository) GetTransactionMidtrans(ID string) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Product").First(&transaction, ID).Error
	return transaction, err
}

// UPDATE TRANSACTION TO DATABASE
func (r *repository) UpdateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Preload("User").Preload("Cart").Preload("Cart.Product").Save(&transaction).Error
	return transaction, err
}

// UPDATE TRANSACTION FROM MIDTRANS TO DATABASE
func (r *repository) UpdateTransactionMidtrans(status string, ID int) error {
	var transaction models.Transaction
	r.db.Preload("Cart.Product").First(&transaction, ID)
	transaction.Status = status
	err := r.db.Debug().Save(&transaction).Error
	return err
}
