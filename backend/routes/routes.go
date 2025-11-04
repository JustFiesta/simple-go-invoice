package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// CORS middleware for frontend
	r.Use(corsMiddleware())

	// API group
	api := r.Group("/api")
	{
		// Health check
		api.GET("/hello", controllers.Ping)

		// Invoice routes
		invoices := api.Group("/invoices")
		{
			invoices.GET("", controllers.GetInvoices)
			invoices.GET("/:id", controllers.GetInvoiceByID)
			invoices.POST("", controllers.CreateInvoice)
			invoices.PUT("/:id", controllers.UpdateInvoice)
			invoices.DELETE("/:id", controllers.DeleteInvoice)
			invoices.GET("/:id/pdf", controllers.GetInvoicePDF)
		}
	}
}

// corsMiddleware handles CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}