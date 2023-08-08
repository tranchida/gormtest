package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	database, err := gorm.Open(mysql.Open("go1:go1@tcp(localhost:3306)/sample?parseTime=true"), &gorm.Config{})
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
		})
	}

	DB = database
}
