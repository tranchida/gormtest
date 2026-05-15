package main

import (
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"tranchida.github.com/gormtest/internal/models"
)

func TestInitDBIsIdempotent(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := initDB(database, true); err != nil {
		t.Fatalf("first initDB: %v", err)
	}

	if err := initDB(database, true); err != nil {
		t.Fatalf("second initDB: %v", err)
	}

	var recettesCount int64
	if err := database.Model(&models.Recette{}).Count(&recettesCount).Error; err != nil {
		t.Fatalf("count recettes: %v", err)
	}
	if recettesCount != 2 {
		t.Fatalf("expected 2 recettes, got %d", recettesCount)
	}

	var livresCount int64
	if err := database.Model(&models.Livre{}).Count(&livresCount).Error; err != nil {
		t.Fatalf("count livres: %v", err)
	}
	if livresCount != 1 {
		t.Fatalf("expected 1 livre, got %d", livresCount)
	}

	var livre models.Livre
	if err := database.Preload("Recettes").Where("titre = ?", "Recettes du poulet").First(&livre).Error; err != nil {
		t.Fatalf("load livre: %v", err)
	}
	if len(livre.Recettes) != 2 {
		t.Fatalf("expected 2 associated recettes, got %d", len(livre.Recettes))
	}
}

func TestInitDBWithoutSeed(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}

	if err := initDB(database, false); err != nil {
		t.Fatalf("initDB without seed: %v", err)
	}

	var recettesCount int64
	if err := database.Model(&models.Recette{}).Count(&recettesCount).Error; err != nil {
		t.Fatalf("count recettes: %v", err)
	}
	if recettesCount != 0 {
		t.Fatalf("expected 0 recettes without seed, got %d", recettesCount)
	}

	var livresCount int64
	if err := database.Model(&models.Livre{}).Count(&livresCount).Error; err != nil {
		t.Fatalf("count livres: %v", err)
	}
	if livresCount != 0 {
		t.Fatalf("expected 0 livres without seed, got %d", livresCount)
	}
}

func TestEnvBool(t *testing.T) {
	t.Setenv("SEED_DB", "false")

	seedEnabled, err := envBool("SEED_DB", true)
	if err != nil {
		t.Fatalf("envBool returned error: %v", err)
	}
	if seedEnabled {
		t.Fatal("expected SEED_DB=false to disable seed")
	}
}
