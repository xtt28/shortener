package database

import (
	"github.com/xtt28/shortener/database/models"
	"gorm.io/gorm"
)

// MigrateAllModels invokes the AutoMigrate method of the gorm.DB on all of the
// application models.
func MigrateAllModels(db *gorm.DB) {
	db.AutoMigrate(&models.ShortLink{})
}
