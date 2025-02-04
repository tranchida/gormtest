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

	s, err := createNewServer()
	if err != nil {
		panic(err)
	}

	if err = s.engine.Run(); err != nil {
		panic(err)
	}

}

func createNewServer() (*server, error) {

	var url string
	if url = os.Getenv("POSTGRESQL_URL"); url == "" {
		url = "postgres://gouser:password@localhost:5432/mydb?sslmode=disable"
	}

	database, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	initDB(database)

	s := &server{
		engine: gin.Default(),
		database: database,
	}

	s.engine.SetTrustedProxies(nil)
	s.engine.Static("/assets", "./assets")
	s.engine.LoadHTMLGlob("templates/*.html")

	s.engine.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	s.engine.GET("/livres", s.AllLivres)
	s.engine.GET("/recettes", s.AllRecettes)
	s.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return s, nil

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

type server struct {
	engine 	  *gin.Engine
	database *gorm.DB
}

func (h server) Index(c *gin.Context) {
	var recettes []models.Recette
	result := h.database.Find(&recettes)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.HTML(http.StatusOK, "index.html", nil)

}

func (h server) AllLivres(c *gin.Context) {

	var livres []models.Livre
	result := h.database.Model(&models.Livre{}).Preload("Recettes").Find(&livres)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.HTML(http.StatusOK, "livresTable", gin.H{
		"data": livres,
	})
}

func (h server) AllRecettes(c *gin.Context) {

	var recettes []models.Recette
	result := h.database.Find(&recettes)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.HTML(http.StatusOK, "recettesTable", gin.H{
		"data": recettes,
	})
}
