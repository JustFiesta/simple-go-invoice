package models

import (
    "time"
    "gorm.io/gorm"
)

type Invoice struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
    
    InvoiceNumber string    `gorm:"uniqueIndex;not null" json:"invoice_number"`
    ClientName    string    `gorm:"not null" json:"client_name"`
    Amount        float64   `gorm:"not null" json:"amount"`
    Status        string    `gorm:"default:'draft'" json:"status"` // draft, sent, paid
    IssueDate     time.Time `json:"issue_date"`
    DueDate       time.Time `json:"due_date"`
}