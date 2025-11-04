package main

import (
	"log"

	"backend/config"
	"backend/database"
	"backend/routes"

	"github.com/gin-gonic/gin"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
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
	routes.SetupRoutes(r, Version, Commit, BuildDate)

	// Start server
	log.Printf("Server starting on %s, version %s", cfg.Port, Version)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}