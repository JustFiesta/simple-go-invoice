package models

import "gorm.io/gorm"

type InvoiceItem struct {
	gorm.Model
	InvoiceID   uint
	Description string
	Quantity    float64
	UnitPrice   float64
	VATRate     float64
}
