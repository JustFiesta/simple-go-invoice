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
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/app/data/invoices.db"
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
		log.Fatal("Invoice Migration failed:", err)
	}

	// Auto-migrate the InvoiceItem model
	err = DB.AutoMigrate(&models.InvoiceItem{})
	if err != nil {
		log.Fatal("InvoiceItem Migration failed:", err)
	}

	log.Println("Database connected and migrated successfully")
}