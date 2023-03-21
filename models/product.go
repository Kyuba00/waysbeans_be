package models

import "time"

type Product struct {
	ID          int       `json:"id" gorm:"primary_key: auto_increment"`
	Name        string    `json:"name" gorm:"type: varchar(50)"`
	Price       int       `json:"price" gorm:"type: int"`
	Image       string    `json:"image" gorm:"type: varchar(255)"`
	Description string    `json:"description" gorm:"type: text"`
	Stock       int       `json:"stock" gorm:"type: int"`
	CreateAt    time.Time `json:"-"`
	UpdateAt    time.Time `json:"-"`
}
