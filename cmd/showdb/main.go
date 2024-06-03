package main

import (
	"fmt"
	"log"

	"tranchida.github.com/gormtest/internal/models"
	"tranchida.github.com/gormtest/internal/services"
)

func main() {

	models.ConnectDatabase()

	products, err := services.FindAll()
	if err != nil {
		log.Fatal(err)

	}
	for _, product := range *products {

		fmt.Println(product.Price)

	}

}
