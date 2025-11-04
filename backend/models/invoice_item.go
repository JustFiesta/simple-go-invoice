package models

import "gorm.io/gorm"

type InvoiceItem struct {
	gorm.Model
	ID          uint    `json:"id" gorm:"primaryKey"`
	InvoiceID   uint    `json:"invoice_id"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	VATRate     float64 `json:"vat_rate"`
}
