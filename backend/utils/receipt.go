package utils

import (
	"os"
	"path/filepath"
	"pg-manager-backend/models"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

func GenerateReceipt(payment models.Payment, tenantName string) (string, error) {
	// P = Portrait, mm = millimeters, A5 size
	pdf := gofpdf.New("P", "mm", "A5", "")
	pdf.AddPage()

	// 1. Branding - Direct string usage
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(24, 144, 255)
	pdf.Cell(0, 10, "PAYMENT RECEIPT")
	pdf.Ln(12)

	// 2. Body - Used strconv for ID conversion (faster than Sprintf)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, "Receipt No: #PAY-"+strconv.FormatUint(uint64(payment.ID), 10))
	pdf.Ln(6)
	pdf.Cell(0, 10, "Date: "+payment.Date.Format("02-Jan-2006"))
	pdf.Ln(10)

	// Table Header
	pdf.SetFillColor(245, 247, 250)
	pdf.CellFormat(0, 10, "Payment Details", "1", 1, "L", true, 0, "")

	// Details
	pdf.Cell(0, 10, "Tenant Name: "+tenantName)
	pdf.Ln(8)
	pdf.Cell(0, 10, "Particulars: Room Rent / Maintenance")
	pdf.Ln(8)

	// Amount
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Total Amount: INR "+strconv.FormatFloat(payment.Amount, 'f', 2, 64))
	pdf.Ln(12)

	// 3. Save Logic
	fileName := "receipt_" + strconv.FormatUint(uint64(payment.ID), 10) + ".pdf"
	dir := filepath.Join("public", "receipts")

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	fullPath := filepath.Join(dir, fileName)
	err := pdf.OutputFileAndClose(fullPath)

	return fileName, err
}
