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

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, fmt.Sprintf("VAT Invoice: %s", inv.InvoiceNumber))
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(100, 8, fmt.Sprintf("Client: %s", inv.ClientName))
	pdf.Ln(8)
	pdf.Cell(100, 8, fmt.Sprintf("Issue Date: %s", inv.IssueDate.Format("2006-01-02")))
	pdf.Ln(8)
	pdf.Cell(100, 8, fmt.Sprintf("Due Date: %s", inv.DueDate.Format("2006-01-02")))
	pdf.Ln(12)

	// Table Headers
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(230, 230, 230)
	cols := []string{"No.", "Description", "Quantity", "Net Price", "VAT %", "Gross Value"}
	widths := []float64{10, 80, 20, 25, 20, 30}
	for i, str := range cols {
		pdf.CellFormat(widths[i], 8, str, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 10)
	var total float64
	for i, item := range inv.Items {
		grossValue := item.Quantity * item.UnitPrice * (1 + item.VATRate/100)
		total += grossValue

		row := []string{
			fmt.Sprintf("%d", i+1),
			item.Description,
			fmt.Sprintf("%.2f", item.Quantity),
			fmt.Sprintf("%.2f", item.UnitPrice),
			fmt.Sprintf("%.0f%%", item.VATRate),
			fmt.Sprintf("%.2f", grossValue),
		}
		for j, str := range row {
			align := "C"
			if j == 1 {
				align = "L"
			}
			pdf.CellFormat(widths[j], 8, str, "1", 0, align, false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Total cost
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(155, 10, "Total Amount Due:", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 10, fmt.Sprintf("%.2f PLN", total), "1", 0, "R", false, 0, "")
	pdf.Ln(15)

	// Footer
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))

	buf := new(bytes.Buffer)
	if err := pdf.Output(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}