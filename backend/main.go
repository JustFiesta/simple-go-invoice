package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "backend/config"
    "backend/database"
    "backend/routes"
)

func main() {
    // Config
    cfg := config.Load()
    
    // DB connection
    database.Connect()
    
    // Router
    r := gin.Default()
    routes.SetupRoutes(r)
    
    // Server start
    log.Printf("Server starting on %s", cfg.Port)
    if err := r.Run(cfg.Port); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
// reapply 4