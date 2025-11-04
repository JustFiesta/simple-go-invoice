package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"backend/database"
	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
)

// Ping - health check endpoint (renamed from /hello to /health)
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "invoice-api",
		"version": "1.0.0",
	})
}

// GetInvoices - GET /invoices?page=1&limit=10&search=&status=&sort=created_at&order=desc
func GetInvoices(c *gin.Context) {
	var invoices []models.Invoice

	// Pagination
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Filters
	search := strings.TrimSpace(c.Query("search"))
	status := strings.TrimSpace(c.Query("status"))

	// Sorting with whitelist
	sort := c.DefaultQuery("sort", "created_at")
	order := strings.ToUpper(c.DefaultQuery("order", "DESC"))

	// Whitelist allowed sort fields
	allowedSortFields := map[string]bool{
		"id":             true,
		"created_at":     true,
		"updated_at":     true,
		"invoice_number": true,
		"client_name":    true,
		"amount":         true,
		"status":         true,
		"issue_date":     true,
		"due_date":       true,
	}

	if !allowedSortFields[sort] {
		sort = "created_at"
	}

	if order != "ASC" && order != "DESC" {
		order = "DESC"
	}

	query := database.DB.Model(&models.Invoice{})

	// Search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("invoice_number LIKE ? OR client_name LIKE ?",
			searchPattern, searchPattern)
	}

	// Status filter
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total records for pagination
	var total int64
	if err := query.Count(&total).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to count invoices",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Query with sorting, limit and offset
	if err := query.Order(sort + " " + order).Limit(limit).Offset(offset).Find(&invoices).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to fetch invoices",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.Header("Cache-Control", "public, max-age=60")

	selfLink := fmt.Sprintf("/api/invoices?page=%d&limit=%d", page, limit)
	response := models.PaginatedResponse(invoices, page, limit, total, selfLink)
	
	c.JSON(http.StatusOK, response)
}

// GetInvoiceByID - GET /invoices/:id
func GetInvoiceByID(c *gin.Context) {
	var invoice models.Invoice
	id := c.Param("id")

	if err := database.DB.First(&invoice, id).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusNotFound,
			"Invoice not found",
			fmt.Sprintf("Invoice with ID %s does not exist", id),
		)
		c.JSON(http.StatusNotFound, response)
		return
	}

	c.Header("Cache-Control", "public, max-age=300")

	related := map[string]string{
		"items": fmt.Sprintf("/api/invoices/%s/items", id),
		"pdf":   fmt.Sprintf("/api/invoices/%s/pdf", id),
	}

	response := models.SingleResourceResponse(
		invoice,
		fmt.Sprintf("/api/invoices/%s", id),
		related,
	)

	c.JSON(http.StatusOK, response)
}

// CreateInvoice - POST /invoices
func CreateInvoice(c *gin.Context) {
	var invoice models.Invoice

	if err := c.ShouldBindJSON(&invoice); err != nil {
		response := models.ErrorResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate required fields
	if err := validateInvoice(&invoice); err != nil {
		response := models.ErrorResponse(
			http.StatusUnprocessableEntity,
			"Validation failed",
			err.Error(),
		)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Check for duplicate invoice number
	var count int64
	database.DB.Model(&models.Invoice{}).Where("invoice_number = ?", invoice.InvoiceNumber).Count(&count)
	if count > 0 {
		response := models.ErrorResponse(
			http.StatusConflict,
			"Invoice number already exists",
			fmt.Sprintf("Invoice number '%s' is already in use", invoice.InvoiceNumber),
		)
		c.JSON(http.StatusConflict, response)
		return
	}

	if err := database.DB.Create(&invoice).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to create invoice",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	location := fmt.Sprintf("/api/invoices/%d", invoice.ID)
	c.Header("Location", location)

	related := map[string]string{
		"items": fmt.Sprintf("/api/invoices/%d/items", invoice.ID),
		"pdf":   fmt.Sprintf("/api/invoices/%d/pdf", invoice.ID),
	}

	response := models.SingleResourceResponse(
		invoice,
		location,
		related,
	)

	c.JSON(http.StatusCreated, response)
}

// UpdateInvoice - PUT /invoices/:id
func UpdateInvoice(c *gin.Context) {
	var invoice models.Invoice
	id := c.Param("id")

	if err := database.DB.First(&invoice, id).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusNotFound,
			"Invoice not found",
			fmt.Sprintf("Invoice with ID %s does not exist", id),
		)
		c.JSON(http.StatusNotFound, response)
		return
	}

	var updateData models.Invoice
	if err := c.ShouldBindJSON(&updateData); err != nil {
		response := models.ErrorResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate updated data
	if err := validateInvoice(&updateData); err != nil {
		response := models.ErrorResponse(
			http.StatusUnprocessableEntity,
			"Validation failed",
			err.Error(),
		)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Check for duplicate invoice number (excluding current invoice)
	if updateData.InvoiceNumber != invoice.InvoiceNumber {
		var count int64
		database.DB.Model(&models.Invoice{}).Where("invoice_number = ? AND id != ?", updateData.InvoiceNumber, id).Count(&count)
		if count > 0 {
			response := models.ErrorResponse(
				http.StatusConflict,
				"Invoice number already exists",
				fmt.Sprintf("Invoice number '%s' is already in use", updateData.InvoiceNumber),
			)
			c.JSON(http.StatusConflict, response)
			return
		}
	}

	// Update fields
	invoice.InvoiceNumber = updateData.InvoiceNumber
	invoice.ClientName = updateData.ClientName
	invoice.Amount = updateData.Amount
	invoice.Status = updateData.Status
	invoice.IssueDate = updateData.IssueDate
	invoice.DueDate = updateData.DueDate

	if err := database.DB.Save(&invoice).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to update invoice",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	related := map[string]string{
		"items": fmt.Sprintf("/api/invoices/%s/items", id),
		"pdf":   fmt.Sprintf("/api/invoices/%s/pdf", id),
	}

	response := models.SingleResourceResponse(
		invoice,
		fmt.Sprintf("/api/invoices/%s", id),
		related,
	)

	c.JSON(http.StatusOK, response)
}

// DeleteInvoice - DELETE /invoices/:id
func DeleteInvoice(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.Invoice{}, id)
	if result.Error != nil {
		response := models.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to delete invoice",
			result.Error.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if result.RowsAffected == 0 {
		response := models.ErrorResponse(
			http.StatusNotFound,
			"Invoice not found",
			fmt.Sprintf("Invoice with ID %s does not exist", id),
		)
		c.JSON(http.StatusNotFound, response)
		return
	}

	c.Status(http.StatusNoContent)
}

// validateInvoice validates invoice data
func validateInvoice(invoice *models.Invoice) error {
	if strings.TrimSpace(invoice.InvoiceNumber) == "" {
		return fmt.Errorf("invoice number is required")
	}

	if strings.TrimSpace(invoice.ClientName) == "" {
		return fmt.Errorf("client name is required")
	}

	if invoice.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	validStatuses := map[string]bool{"draft": true, "sent": true, "paid": true}
	if !validStatuses[invoice.Status] {
		return fmt.Errorf("status must be one of: draft, sent, paid")
	}

	if invoice.DueDate.Before(invoice.IssueDate) {
		return fmt.Errorf("due date cannot be before issue date")
	}

	return nil
}

// GetInvoicePDF - GET /invoices/:id/pdf
func GetInvoicePDF(c *gin.Context) {
	var invoice models.Invoice
	id := c.Param("id")

	if err := database.DB.Preload("Items").First(&invoice, id).Error; err != nil {
		response := models.ErrorResponse(
			http.StatusNotFound,
			"Invoice not found",
			fmt.Sprintf("Invoice with ID %s does not exist", id),
		)
		c.JSON(http.StatusNotFound, response)
		return
	}

	pdfBytes, err := services.GenerateInvoicePDF(invoice)
	if err != nil {
		response := models.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to generate PDF",
			err.Error(),
		)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"invoice-%s.pdf\"", invoice.InvoiceNumber))
	c.Header("Cache-Control", "private, max-age=3600")
	
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}