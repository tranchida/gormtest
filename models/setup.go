package models

import (
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

	if p.ID == 0 {
		database.Create([]Product{
			{Code: "sample1", Price: 10},
			{Code: "sample2", Price: 20},
			{Code: "sample3", Price: 30},
			{Code: "sample4", Price: 40},
			{Code: "sample5", Price: 50},
		})
	}

	DB = database
}
