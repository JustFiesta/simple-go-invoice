package controllers

import (
	"net/http"
	"backend/database"
	"backend/models"
	"github.com/gin-gonic/gin"
)

// AddInvoiceItem - POST /invoices/:id/items
func AddInvoiceItem(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice

	if err := database.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	var item models.InvoiceItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.InvoiceID = invoice.ID

	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		return
	}

	updateInvoiceTotal(&invoice)

	c.JSON(http.StatusCreated, item)
}

// GetInvoiceItems - GET /invoices/:id/items
func GetInvoiceItems(c *gin.Context) {
	id := c.Param("id")
	var items []models.InvoiceItem

	if err := database.DB.Where("invoice_id = ?", id).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}

	c.JSON(http.StatusOK, items)
}

// UpdateInvoiceItem - PUT /invoices/:id/items/:itemId
func UpdateInvoiceItem(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")

	var invoice models.Invoice
	if err := database.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	var item models.InvoiceItem
	if err := database.DB.First(&item, itemId).Error; err != nil || item.InvoiceID != invoice.ID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	var updated models.InvoiceItem
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.Description = updated.Description
	item.Quantity = updated.Quantity
	item.UnitPrice = updated.UnitPrice
	item.VATRate = updated.VATRate

	if err := database.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	updateInvoiceTotal(&invoice)
	c.JSON(http.StatusOK, item)
}

// DeleteInvoiceItem - DELETE /invoices/:id/items/:itemId
func DeleteInvoiceItem(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")

	var invoice models.Invoice
	if err := database.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	if err := database.DB.Delete(&models.InvoiceItem{}, itemId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	updateInvoiceTotal(&invoice)
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}

// Helper function to update invoice total
func updateInvoiceTotal(invoice *models.Invoice) {
	var total float64
	database.DB.Model(&models.InvoiceItem{}).
		Where("invoice_id = ?", invoice.ID).
		Select("SUM(quantity * unit_price)").Scan(&total)

	invoice.Amount = total
	database.DB.Save(&invoice)
}