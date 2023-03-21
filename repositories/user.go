package repositories

import (
	"waysbeans_be/models"

	"gorm.io/gorm"
)

// USER REPOSITORY INTERFACE
type UserRepository interface {
	GetUser(ID int) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
}

// USER REPOSITORY FUNCTION
func RepositoryUser(db *gorm.DB) *repository {
	return &repository{db}
}

// FIND USER FROM DATABASE
func (r *repository) GetUser(ID int) (models.User, error) {
	var user models.User
	err := r.db.First(&user, ID).Error
	return user, err
}

// UPDATE USER TO DATABASE
func (r *repository) UpdateUser(user models.User) (models.User, error) {
	err := r.db.Save(&user).Error
	return user, err
}
