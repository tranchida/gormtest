package main

import (
	"tranchida.github.com/gormtest/controlers"
	"tranchida.github.com/gormtest/models"
)

func main() {

	models.ConnectDatabase()

	controlers.SetupRouter()
}
