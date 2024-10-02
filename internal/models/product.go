package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string `json:"code" gorm:"size:150"`
	Price uint   `json:"price"`
	Stock uint   `json:"stock"`
	Image []byte `json:"image"`
}
