package services

import (
	"bytes"
	"fmt"
	"time"

	"github.com/phpdave11/gofpdf"
	"backend/models"
)

func GenerateInvoicePDF(inv models.Invoice) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	// Header
	pdf.Cell(190, 10, "VAT Invoice")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Invoice number: %s", inv.InvoiceNumber))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Client: %s", inv.ClientName))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Total: %.2f PLN", inv.Amount))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Status: %s", inv.Status))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Issued at: %s", inv.IssueDate.Format("2006-01-02")))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Payment due date: %s", inv.DueDate.Format("2006-01-02")))
	pdf.Ln(12)

	// Footer
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated at: %s", time.Now().Format("2006-01-02 15:04:05")))

	// Return PDF as byte slice
	buf := new(bytes.Buffer)
	if err := pdf.Output(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}