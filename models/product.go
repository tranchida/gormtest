package models

import "time"

type OwnModel struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-";sql:"index"`
}

type Product struct {
	OwnModel
	Code  string `json:"code" gorm:"size:150"`
	Price uint   `json:"price"`
}
