package routes

import (
    "github.com/gin-gonic/gin"
    "backend/controllers"
)

func SetupRoutes(r *gin.Engine) {
    // CORS for Vue
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    })
    
    api := r.Group("/api")
    {
		api.GET("/hello", controllers.Ping)
        invoices := api.Group("/invoices")
        {
            invoices.GET("", controllers.GetInvoices)
            invoices.GET("/:id", controllers.GetInvoiceByID)
            invoices.POST("", controllers.CreateInvoice)
            invoices.PUT("/:id", controllers.UpdateInvoice)
            invoices.DELETE("/:id", controllers.DeleteInvoice)
        }
    }
}