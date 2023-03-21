package models

import "time"

type User struct {
	ID       int       `json:"id" gorm:"primary_key: auto_increment"`
	Name     string    `json:"name" gorm:"type: varchar(50)"`
	Email    string    `json:"email" gorm:"type: varchar(50);unique"`
	Password string    `json:"password" gorm:"type: varchar(100)"`
	Image    string    `json:"image" gorm:"type: varchar(255)"`
	Phone    string    `json:"phone" gorm:"type: varchar(50)"`
	Address  string    `json:"address" gorm:"type: text"`
	Role     string    `json:"role" gorm:"type: varchar(50)"`
	CreateAt time.Time `json:"-"`
	UpdateAt time.Time `json:"-"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (UserResponse) TableName() string {
	return "users"
}
