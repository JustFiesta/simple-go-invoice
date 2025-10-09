package database

import (
    "log"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "backend/models"
)

var DB *gorm.DB

func Connect() {
    var err error
    DB, err = gorm.Open(sqlite.Open("invoices.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

	// Auto-migrate the Invoice model
    err = DB.AutoMigrate(&models.Invoice{})
    if err != nil {
        log.Fatal("Migration failed:", err)
    }

    log.Println("Database connected and migrated")
}