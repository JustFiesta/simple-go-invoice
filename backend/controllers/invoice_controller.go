package controllers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "backend/database"
    "backend/models"
)

func Ping(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// GET /invoices?page=1&limit=10&search=&status=&sort=created_at&order=desc
func GetInvoices(c *gin.Context) {
    var invoices []models.Invoice
    
    // Pagination
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset := (page - 1) * limit
    
    // Filters
    search := c.Query("search")
    status := c.Query("status")
    
    // Sorting
    sort := c.DefaultQuery("sort", "created_at")
    order := c.DefaultQuery("order", "desc")
    
    query := database.DB.Model(&models.Invoice{})
    
    // Search filter
    if search != "" {
        query = query.Where("invoice_number LIKE ? OR client_name LIKE ?", 
            "%"+search+"%", "%"+search+"%")
    }
    
    // Status filter
    if status != "" {
        query = query.Where("status = ?", status)
    }
    
    // Sorting
    query = query.Order(sort + " " + order)
    
    // Count total records for pagination
    var total int64
    query.Count(&total)
    
    // Query with limit and offset
    query.Limit(limit).Offset(offset).Find(&invoices)
    
    c.JSON(http.StatusOK, gin.H{
        "data":  invoices,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}

// GET /invoices/:id
func GetInvoiceByID(c *gin.Context) {
    var invoice models.Invoice
    id := c.Param("id")
    
    if err := database.DB.First(&invoice, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
        return
    }
    
    c.JSON(http.StatusOK, invoice)
}

// POST /invoices
func CreateInvoice(c *gin.Context) {
    var invoice models.Invoice
    
    if err := c.ShouldBindJSON(&invoice); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := database.DB.Create(&invoice).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
        return
    }
    
    c.JSON(http.StatusCreated, invoice)
}

// PUT /invoices/:id
func UpdateInvoice(c *gin.Context) {
    var invoice models.Invoice
    id := c.Param("id")
    
    if err := database.DB.First(&invoice, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
        return
    }
    
    if err := c.ShouldBindJSON(&invoice); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    database.DB.Save(&invoice)
    c.JSON(http.StatusOK, invoice)
}

// DELETE /invoices/:id
func DeleteInvoice(c *gin.Context) {
    id := c.Param("id")
    
    if err := database.DB.Delete(&models.Invoice{}, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted"})
}