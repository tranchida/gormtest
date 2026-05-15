package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"tranchida.github.com/gormtest/internal/models"
)

type appConfig struct {
	port        string
	sqlitePath  string
	seedEnabled bool
}

func main() {

	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	s, err := createNewServer(config)
	if err != nil {
		panic(err)
	}

	if err = s.engine.Run(":" + config.port); err != nil {
		panic(err)
	}

}

func loadConfig() (appConfig, error) {

	if err := loadDotEnv(".env"); err != nil {
		return appConfig{}, err
	}

	seedEnabled, err := envBool("SEED_DB", true)
	if err != nil {
		return appConfig{}, err
	}

	return appConfig{
		port:        envOrDefault("APP_PORT", "8080"),
		sqlitePath:  envOrDefault("SQLITE_PATH", "gorm.db"),
		seedEnabled: seedEnabled,
	}, nil
}

func loadDotEnv(path string) error {

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	return godotenv.Load(path)
}

func envBool(name string, defaultValue bool) (bool, error) {

	value, ok := os.LookupEnv(name)
	if !ok || value == "" {
		return defaultValue, nil
	}

	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("parse %s: %w", name, err)
	}

	return parsedValue, nil
}

func envOrDefault(name, defaultValue string) string {

	if value, ok := os.LookupEnv(name); ok && value != "" {
		return value
	}

	return defaultValue
}

func createNewServer(config appConfig) (*server, error) {

	databaseDir := filepath.Dir(config.sqlitePath)
	if databaseDir != "." {
		if err := os.MkdirAll(databaseDir, 0o755); err != nil {
			return nil, err
		}
	}

	database, err := gorm.Open(sqlite.Open(config.sqlitePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}

	if err = initDB(database, config.seedEnabled); err != nil {
		return nil, err
	}

	s := &server{
		engine:   gin.Default(),
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

func initDB(database *gorm.DB, seedEnabled bool) error {

	if err := database.AutoMigrate(&models.Livre{}, &models.Recette{}); err != nil {
		return err
	}

	if !seedEnabled {
		return nil
	}

	return seedDB(database)
}

func seedDB(database *gorm.DB) error {

	type recetteSeed struct {
		Nom    string
		Niveau uint
		Temps  uint
	}

	recettesSeed := []recetteSeed{
		{
			Nom:    "Poulet au curry",
			Niveau: 2,
			Temps:  30,
		},
		{
			Nom:    "Poulet au citron",
			Niveau: 1,
			Temps:  20,
		},
	}

	return database.Transaction(func(tx *gorm.DB) error {
		recettes := make([]models.Recette, 0, len(recettesSeed))

		for _, recetteSeed := range recettesSeed {
			recette := models.Recette{}
			err := tx.Where("nom = ?", recetteSeed.Nom).First(&recette).Error
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}

				recette = models.Recette{
					Nom:    recetteSeed.Nom,
					Niveau: recetteSeed.Niveau,
					Temps:  recetteSeed.Temps,
				}
				if err := tx.Create(&recette).Error; err != nil {
					return err
				}
			} else {
				recette.Niveau = recetteSeed.Niveau
				recette.Temps = recetteSeed.Temps
				if err := tx.Save(&recette).Error; err != nil {
					return err
				}
			}

			recettes = append(recettes, recette)
		}

		livre := models.Livre{}
		err := tx.Where("titre = ?", "Recettes du poulet").First(&livre).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			livre = models.Livre{
				Titre: "Recettes du poulet",
			}
			if err := tx.Create(&livre).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&livre).Association("Recettes").Replace(recettes); err != nil {
			return err
		}

		return nil
	})
}

type server struct {
	engine   *gin.Engine
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
