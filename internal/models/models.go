package models

import (
	"gorm.io/gorm"
)

type Livre struct {
	gorm.Model
	Titre    string    `json:"titre" gorm:"size:150"`
	Recettes []Recette `gorm:"many2many:livre_recette;"`
}

type Recette struct {
	gorm.Model
	Nom    string `json:"nom" gorm:"size:150"`
	Niveau uint   `json:"niveau"`
	Temps  uint   `json:"temps"`
}
