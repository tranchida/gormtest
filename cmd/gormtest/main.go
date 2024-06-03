package main

import (
	"tranchida.github.com/gormtest/internal/controlers"
	"tranchida.github.com/gormtest/internal/models"
)

func main() {

	models.ConnectDatabase()

	controlers.SetupRouter()
}
