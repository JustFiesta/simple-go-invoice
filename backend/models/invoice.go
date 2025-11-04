package models

import (
    "time"
    "gorm.io/gorm"
)

type Invoice struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	InvoiceNumber string         `gorm:"uniqueIndex;not null" json:"invoice_number" binding:"required"`
	ClientName    string         `gorm:"not null" json:"client_name" binding:"required"`
	Amount        float64        `gorm:"not null" json:"amount" binding:"required,gt=0"`
	Status        string         `gorm:"default:'draft'" json:"status" binding:"required,oneof=draft sent paid"`
	IssueDate     time.Time      `json:"issue_date" binding:"required"`
	DueDate       time.Time      `json:"due_date" binding:"required"`
}