package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"tranchida.github.com/gormtest/internal/models"
)

func main() {

	var url string
	if url = os.Getenv("POSTGRESQL_URL"); url == "" {
		url = "postgres://gouser:password@localhost:5432/mydb?sslmode=disable"
	}

	database, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic("failed to connect database")
	}

	initDB(database)

	handler := handler{database: database}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/livres", handler.AllLivres)
	r.GET("/recettes", handler.AllRecettes)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run()

}

func initDB(database *gorm.DB) {

	database.AutoMigrate(&models.Livre{}, &models.Recette{})

	recette := models.Recette{}
	result := database.First(&recette)
	if result.RowsAffected == 0 {

		p1 := models.Recette{
			Nom:    "Poulet au curry",
			Niveau: 2,
			Temps:  30,
		}

		p2 := models.Recette{
			Nom:    "Poulet au citron",
			Niveau: 1,
			Temps:  20,
		}

		l := models.Livre{
			Titre: "Recettes du poulet",
			Recettes: []models.Recette{
				p1,
				p2,
			},
		}

		database.Create(&l)

		database.Save(&l)

	}
}

type handler struct {
	database *gorm.DB
}

func (h handler) AllLivres(c *gin.Context) {

	var livres []models.Livre
	result := h.database.Model(&models.Livre{}).Preload("Recettes").Find(&livres)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, livres)
}

func (h handler) AllRecettes(c *gin.Context) {

	var recettes []models.Recette
	result := h.database.Find(&recettes)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, recettes)
}
