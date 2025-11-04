package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"backend/database"
	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
)

// Ping - health check endpoint
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count invoices"})
		return
	}

	// Query with sorting, limit and offset
	if err := query.Order(sort + " " + order).Limit(limit).Offset(offset).Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  invoices,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetInvoiceByID - GET /invoices/:id
func GetInvoiceByID(c *gin.Context) {
	var invoice models.Invoice
	id := c.Param("id")

	if err := database.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// CreateInvoice - POST /invoices
func CreateInvoice(c *gin.Context) {
	var invoice models.Invoice

	if err := c.ShouldBindJSON(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if err := validateInvoice(&invoice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate invoice number
	var count int64
	database.DB.Model(&models.Invoice{}).Where("invoice_number = ?", invoice.InvoiceNumber).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Invoice number already exists"})
		return
	}

	if err := database.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}

	c.JSON(http.StatusCreated, invoice)
}

// UpdateInvoice - PUT /invoices/:id
func UpdateInvoice(c *gin.Context) {
	var invoice models.Invoice
	id := c.Param("id")

	if err := database.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	var updateData models.Invoice
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate updated data
	if err := validateInvoice(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate invoice number (excluding current invoice)
	if updateData.InvoiceNumber != invoice.InvoiceNumber {
		var count int64
		database.DB.Model(&models.Invoice{}).Where("invoice_number = ? AND id != ?", updateData.InvoiceNumber, id).Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Invoice number already exists"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
		return
	}

	c.JSON(http.StatusOK, invoice)
}

// DeleteInvoice - DELETE /invoices/:id
func DeleteInvoice(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.Invoice{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invoice"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}

// validateInvoice validates invoice data
func validateInvoice(invoice *models.Invoice) error {
	if strings.TrimSpace(invoice.InvoiceNumber) == "" {
		return gin.Error{Err: nil, Type: gin.ErrorTypePublic, Meta: "Invoice number is required"}
	}

	if strings.TrimSpace(invoice.ClientName) == "" {
		return gin.Error{Err: nil, Type: gin.ErrorTypePublic, Meta: "Client name is required"}
	}

	if invoice.Amount <= 0 {
		return gin.Error{Err: nil, Type: gin.ErrorTypePublic, Meta: "Amount must be greater than 0"}
	}

	validStatuses := map[string]bool{"draft": true, "sent": true, "paid": true}
	if !validStatuses[invoice.Status] {
		return gin.Error{Err: nil, Type: gin.ErrorTypePublic, Meta: "Status must be one of: draft, sent, paid"}
	}

	if invoice.DueDate.Before(invoice.IssueDate) {
		return gin.Error{Err: nil, Type: gin.ErrorTypePublic, Meta: "Due date cannot be before issue date"}
	}

	return nil
}

// GetInvoicePDF - GET /invoices/:id/pdf
func GetInvoicePDF(c *gin.Context) {
    var invoice models.Invoice
    id := c.Param("id")

	if err := database.DB.Preload("Items").First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

    pdfBytes, err := services.GenerateInvoicePDF(invoice)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
        return
    }

    c.Data(http.StatusOK, "application/pdf", pdfBytes)
}