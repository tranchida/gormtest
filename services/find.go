package services

import (
	"tranchida.github.com/gormtest/models"
)

func FindAll() (*[]models.Product, error) {

	var products []models.Product
	models.DB.Find(&products)
	if models.DB.Error != nil {
		return nil, models.DB.Error
	}

	return &products, nil

}
