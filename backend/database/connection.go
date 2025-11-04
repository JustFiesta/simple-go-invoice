package database

import (
	"log"
	"os"

	"backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error
	dbPath := "/app/data/invoices.db"

	// Ensure data directory exists
	dataDir := "/app/data"
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			log.Fatal("Failed to create data directory:", err)
		}
	}

	// Configure GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the Invoice model
	err = DB.AutoMigrate(&models.Invoice{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database connected and migrated successfully")
}