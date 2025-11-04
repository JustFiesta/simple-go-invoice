package main

import (
	"log"

	"backend/config"
	"backend/database"
	"backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Load configuration
	cfg := config.Load()

	// Connect to database
	database.Connect()

	// Initialize router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	log.Printf("Server starting on %s", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}