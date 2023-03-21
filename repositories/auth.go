package repositories

import (
	"waysbeans_be/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(user models.User) (models.User, error)
	Login(email string) (models.User, error)
	GetUserAuth(ID int) (models.User, error)
}

// AUTH REPOSITORY FUNCTION
func RepositoryAuth(db *gorm.DB) *repository {
	return &repository{db}
}

// CREATE USER TO DATABASE
func (r *repository) Register(user models.User) (models.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

// GET USER BY EMAIL FOR LOGIN
func (r *repository) Login(email string) (models.User, error) {
	var user models.User
	err := r.db.First(&user, "email=?", email).Error
	return user, err
}

// GET USER BY ID FOR CHECK AUTH
func (r *repository) GetUserAuth(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error
	return user, err
}
