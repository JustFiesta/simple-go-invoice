package controllers

import (
	"fmt"
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
		response := models.ErrorResponse(
			http.StatusNotFound,
			"Invoice not found",
			fmt.Sprintf("Invoice with ID %s does not exist", id),
		)
		c.JSON(http.StatusNotFound, response)
		return
	}

	var item models.InvoiceItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response := models.ErrorResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	item.InvoiceID = invoice.ID

	if err := database.DB.Create(&item).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to add item",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	updateInvoiceTotal(&invoice)

	location := fmt.Sprintf("/api/invoices/%s/items/%d", id, item.ID)
	c.Header("Location", location)

	related := map[string]string{
		"invoice": fmt.Sprintf("/api/invoices/%s", id),
	}

	response := models.SingleResourceResponse(
		item,
		location,
		related,
	)

	c.JSON(http.StatusCreated, response)
}

// GetInvoiceItems - GET /invoices/:id/items
func GetInvoiceItems(c *gin.Context) {
	id := c.Param("id")
	var items []models.InvoiceItem

	// Check if invoice exists
	var invoice models.Invoice
	if err := database.DB.First(&invoice, id).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusNotFound,
			"Invoice not found",
			fmt.Sprintf("Invoice with ID %s does not exist", id),
		)
		c.JSON(http.StatusNotFound, response)
		return
	}

	if err := database.DB.Where("invoice_id = ?", id).Find(&items).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to fetch items",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.Header("Cache-Control", "public, max-age=60")

	related := map[string]string{
		"invoice": fmt.Sprintf("/api/invoices/%s", id),
	}

	response := models.SingleResourceResponse(
		items,
		fmt.Sprintf("/api/invoices/%s/items", id),
		related,
	)

	c.JSON(http.StatusOK, response)
}

// UpdateInvoiceItem - PUT /invoices/:id/items/:itemId
func UpdateInvoiceItem(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")

	var invoice models.Invoice
	if err := database.DB.First(&invoice, id).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusNotFound,
			"Invoice not found",
			fmt.Sprintf("Invoice with ID %s does not exist", id),
		)
		c.JSON(http.StatusNotFound, response)
		return
	}

	var item models.InvoiceItem
	if err := database.DB.First(&item, itemId).Error; err != nil || item.InvoiceID != invoice.ID {
		response := models.ErrorResponse(
			http.StatusNotFound,
			"Item not found",
			fmt.Sprintf("Item with ID %s not found for invoice %s", itemId, id),
		)
		c.JSON(http.StatusNotFound, response)
		return
	}

	var updated models.InvoiceItem
	if err := c.ShouldBindJSON(&updated); err != nil {
		response := models.ErrorResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	item.Description = updated.Description
	item.Quantity = updated.Quantity
	item.UnitPrice = updated.UnitPrice
	item.VATRate = updated.VATRate

	if err := database.DB.Save(&item).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to update item",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	updateInvoiceTotal(&invoice)

	related := map[string]string{
		"invoice": fmt.Sprintf("/api/invoices/%s", id),
		"all_items": fmt.Sprintf("/api/invoices/%s/items", id),
	}

	response := models.SingleResourceResponse(
		item,
		fmt.Sprintf("/api/invoices/%s/items/%s", id, itemId),
		related,
	)

	c.JSON(http.StatusOK, response)
}

// DeleteInvoiceItem - DELETE /invoices/:id/items/:itemId
func DeleteInvoiceItem(c *gin.Context) {
	id := c.Param("id")
	itemId := c.Param("itemId")

	var invoice models.Invoice
	if err := database.DB.First(&invoice, id).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusNotFound,
			"Invoice not found",
			fmt.Sprintf("Invoice with ID %s does not exist", id),
		)
		c.JSON(http.StatusNotFound, response)
		return
	}

	result := database.DB.Delete(&models.InvoiceItem{}, itemId)
	if result.Error != nil {
		response := models.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to delete item",
			result.Error.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if result.RowsAffected == 0 {
		response := models.ErrorResponse(
			http.StatusNotFound,
			"Item not found",
			fmt.Sprintf("Item with ID %s does not exist", itemId),
		)
		c.JSON(http.StatusNotFound, response)
		return
	}

	updateInvoiceTotal(&invoice)

	c.Status(http.StatusNoContent)
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