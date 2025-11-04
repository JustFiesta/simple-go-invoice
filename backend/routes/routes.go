package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, version, commit, buildDate string) {
	// CORS middleware for frontend
	r.Use(corsMiddleware())

	r.Use(cacheMiddleware())

	// API group
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			controllers.Ping(c, version, commit, buildDate)
		})

		// Invoice routes
		invoices := api.Group("/invoices")
		{
			invoices.GET("", controllers.GetInvoices)
			invoices.GET("/:id", controllers.GetInvoiceByID)
			invoices.POST("", controllers.CreateInvoice)
			invoices.PUT("/:id", controllers.UpdateInvoice)
			invoices.DELETE("/:id", controllers.DeleteInvoice)
			invoices.GET("/:id/pdf", controllers.GetInvoicePDF)

			// Invoice items (nested resource)
			invoices.POST("/:id/items", controllers.AddInvoiceItem)
			invoices.GET("/:id/items", controllers.GetInvoiceItems)
			invoices.PUT("/:id/items/:itemId", controllers.UpdateInvoiceItem)
			invoices.DELETE("/:id/items/:itemId", controllers.DeleteInvoiceItem)
		}
	}
}

// corsMiddleware handles CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Location")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func cacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only set cache for GET requests if not already set
		if c.Request.Method == "GET" {
			c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Writer.Header().Set("Pragma", "no-cache")
			c.Writer.Header().Set("Expires", "0")
		}
		
		c.Next()
	}
}