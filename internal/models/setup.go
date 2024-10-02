package models

import (
	"log"
	"math/rand"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	database, err := gorm.Open(sqlite.Open("gormtest.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	database.AutoMigrate(&Product{})

	p := Product{}
	database.First(&p)

	file, err := os.ReadFile("sample.jpg")
	if err != nil {
		log.Fatal(err)
	}

	if p.ID == 0 {
		database.Create([]Product{
			{Code: "sample1", Price: 10, Stock: uint(rand.Intn(100)), Image: file},
			{Code: "sample2", Price: 20, Stock: uint(rand.Intn(100)), Image: file},
			{Code: "sample3", Price: 30, Stock: uint(rand.Intn(100)), Image: file},
			{Code: "sample4", Price: 40, Stock: uint(rand.Intn(100)), Image: file},
			{Code: "sample5", Price: 50, Stock: uint(rand.Intn(100)), Image: file},
		})
	}

	DB = database
}
