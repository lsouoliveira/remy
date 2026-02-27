package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"remy/internal/models"
)

func SetupDatabase(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Note{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
